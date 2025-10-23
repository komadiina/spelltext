package server

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	pbInventory "github.com/komadiina/spelltext/proto/inventory"
	pbRepo "github.com/komadiina/spelltext/proto/repo"
	pb "github.com/komadiina/spelltext/proto/store"
	"github.com/komadiina/spelltext/server/store/config"
	"github.com/komadiina/spelltext/utils/singleton/logging"
	"google.golang.org/grpc"
)

type Connections struct {
	Inventory *grpc.ClientConn
}

type Clients struct {
	Inventory pbInventory.InventoryClient
}

type StoreService struct {
	pb.UnimplementedStoreServer
	DbPool      *pgxpool.Pool
	Config      *config.Config
	Logger      *logging.Logger
	Clients     *Clients
	Connections *Connections
}

func (s *StoreService) ListVendors(ctx context.Context, req *pb.StoreListVendorRequest) (*pb.StoreListVendorResponse, error) {
	sql, _, err := sq.Select("*").From("vendors").ToSql()
	if err != nil {
		s.Logger.Error("failed to build query", "err", err)
		return nil, err
	}

	rows, err := s.DbPool.Query(ctx, sql)
	if err != nil {
		s.Logger.Error("failed to run query", "err", err)
		return nil, err
	}

	var vendors []*pbRepo.Vendor
	for rows.Next() {
		v := &pbRepo.Vendor{}
		err := rows.Scan(&v.Id, &v.Name, &v.WareShorthand)

		if err != nil {
			s.Logger.Error("failed to scan", "err", err)
			return nil, err
		}

		vendors = append(vendors, v)
	}

	return &pb.StoreListVendorResponse{Vendors: vendors}, nil
}

func (s *StoreService) ListVendorItems(ctx context.Context, req *pb.StoreListVendorItemRequest) (*pb.ListVendorItemResponse, error) {
	cte := sq.Select("v.id AS id, vw.item_type_id AS item_type_id, it.code AS code").
		From("vendors AS v").
		InnerJoin("vendor_wares AS vw ON vw.vendor_id = v.id").
		InnerJoin("item_types AS it ON it.id = vw.item_type_id").
		Where("v.id = $1")

	cteSql, _, err := cte.ToSql()
	if err != nil {
		s.Logger.Error("failed to build cte", "err", err)
		return nil, err
	}
	prefix := fmt.Sprintf("WITH vendors_filtered AS (%s)", cteSql)

	query := sq.
		Select("i.*, templ.*, it.*, es.*").
		From("item_templates AS templ").
		InnerJoin("vendors_filtered ON vendors_filtered.item_type_id = templ.item_type_id").
		InnerJoin("item_types as it on it.id = templ.item_type_id").
		InnerJoin("equip_slots as es on es.id = templ.equip_slot_id").
		LeftJoin("items AS i ON i.item_template_id = templ.id")

	sql, _, err := query.ToSql()
	if err != nil {
		s.Logger.Error("failed to build query", "err", err)
		return nil, err
	}

	sql = fmt.Sprintf("%s %s", prefix, sql)
	rows, err := s.DbPool.Query(ctx, sql, req.VendorId)
	if err != nil {
		s.Logger.Error("failed to query", "err", err)
		return nil, err
	}

	var items []*pbRepo.Item
	for rows.Next() {
		var foo *any
		i := &pbRepo.Item{}
		templ := &pbRepo.ItemTemplate{}
		es := &pbRepo.EquipSlot{}
		it := &pbRepo.ItemType{}

		err := rows.Scan(
			&i.Id,
			&i.Prefix,
			&i.Suffix,
			&i.ItemTemplateId,
			&i.Health,
			&i.Power,
			&i.Strength,
			&i.Spellpower,
			&i.BonusDamage,
			&i.BonusArmor,
			&templ.Id,
			&templ.Name,
			&templ.ItemTypeId,
			&templ.EquipSlotId,
			&templ.Description,
			&templ.GoldPrice,
			&templ.BuyableWithTokens,
			&templ.TokenPrice,
			&foo, // metadata, unnecessary atm
			&es.Id,
			&es.Code,
			&es.Name,
			&it.Id,
			&it.Code,
			&it.Name,
		)

		if err != nil {
			s.Logger.Error("failed to scan", "err", err)
			return nil, err
		}

		templ.EquipSlot = es
		templ.ItemType = it
		i.ItemTemplate = templ

		items = append(items, i)
	}

	return &pb.ListVendorItemResponse{Items: items, TotalCount: -1}, nil
}

func (s *StoreService) BuyItems(ctx context.Context, req *pb.BuyItemRequest) (*pb.BuyItemResponse, error) {
	s.Logger.Info(s.DbPool.Config())
	s.Logger.Debugf("BuyItems(%+v)", req)
	start := time.Now()

	errResp := &pb.BuyItemResponse{Success: false, Message: "error ocurred"}
	// check if character has enough gold
	sql, _, err := sq.Select("c.gold").
		From("characters AS c").
		Where("c.character_id = $1").
		ToSql()

	c := &pbRepo.Character{CharacterId: req.CharacterId}
	rows, err := s.DbPool.Query(ctx, sql, req.CharacterId)
	for rows.Next() {
		err := rows.Scan(&c.Gold)

		if err != nil {
			s.Logger.Error(err)
			return errResp, err
		}
	}
	rows.Close()

	// get item gold prices
	sql, args, err := sq.
		Select("i.id, it.gold_price, i.item_template_id").
		From("items AS i").
		InnerJoin("item_templates AS it ON i.id = it.id").
		Where(sq.Eq{"i.id": req.ItemIds}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		s.Logger.Error(err)
		return errResp, err
	}

	rows, err = s.DbPool.Query(ctx, sql, args...)
	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	var items []*pbRepo.Item
	for rows.Next() {
		item := &pbRepo.Item{}

		err := rows.Scan(&item.Id, &item.ItemTemplate.GoldPrice, &item.ItemTemplateId)
		if err != nil {
			s.Logger.Error(err)
			return errResp, err
		}

		items = append(items, item)
	}

	var sum uint64 = 0
	for _, item := range items {
		sum += item.ItemTemplate.GetGoldPrice()
	}

	if sum > c.Gold {
		s.Logger.Infof("character %s overbought attempt, gold_amount=%d, character_gold=%d", c.CharacterName, sum, c.Gold)
		return &pb.BuyItemResponse{Success: false, Message: "not enough gold"}, fmt.Errorf("overbuy attempt")
	}

	batch := &pgx.Batch{}
	sql = "INSERT INTO item_instances VALUES (DEFAULT, $1, $2, DEFAULT, DEFAULT) RETURNING item_instance_id"
	var itemInstanceIds []uint64
	for _, item := range items {
		batch.Queue(sql, item.ItemTemplateId, c.CharacterId)
	}

	res := s.DbPool.SendBatch(ctx, batch)
	defer res.Close()
	for i := 0; i < len(items); i++ {
		row := res.QueryRow()
		if err != nil {
			s.Logger.Error(err)
			return nil, err
		}

		var itemInstanceId uint64
		err = row.Scan(&itemInstanceId)
		if err != nil {
			s.Logger.Error(err)
			return nil, err
		}

		itemInstanceIds = append(itemInstanceIds, itemInstanceId)
	}

	// update character gold
	sql = "UPDATE characters SET gold = gold - $1 WHERE character_id = $2"
	_, err = s.DbPool.Exec(ctx, sql, sum, c.CharacterId)
	if err != nil {
		s.Logger.Error(err)
		return errResp, err
	}

	// TODO: move from direct service-service to MQ (problem: wait-for-ack)
	_, err = s.Clients.Inventory.AddItemsToBackpack(ctx, &pbInventory.AddItemsToBackpackRequest{
		CharacterId:     c.CharacterId,
		ItemInstanceIds: itemInstanceIds,
	})

	if err != nil {
		s.Logger.Error(err)
		return errResp, err
	}

	s.Logger.Infof("finished, start=%v, t=%v", start.Format(time.RFC3339), time.Since(start))

	return &pb.BuyItemResponse{Success: true, Message: "items bought & added to inventory"}, nil
}

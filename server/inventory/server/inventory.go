package server

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	pb "github.com/komadiina/spelltext/proto/inventory"
	pbRepo "github.com/komadiina/spelltext/proto/repo"
	"github.com/komadiina/spelltext/server/inventory/config"
	"github.com/komadiina/spelltext/utils/singleton/logging"
)

type InventoryService struct {
	pb.UnimplementedInventoryServer
	DbPool *pgxpool.Pool
	Config *config.Config
	Logger *logging.Logger
}

func tryConnect(s *InventoryService, context context.Context, conninfo string, backoff time.Duration, maxRetries int, boFormula func(time.Duration) time.Duration) (pgx.Conn, error) {
	try := 1
	for {
		conn, err := pgx.Connect(context, conninfo)

		if err != nil && try >= maxRetries {
			// conn not established, max retries exceeded
			s.Logger.Fatal(err)
		} else if err == nil && try < maxRetries {
			// conn established within maxRetries
			s.Logger.Info("pgpool connection established")
			return *conn, nil
		} else if err != nil && try < maxRetries {
			// conn not established, backoff
			s.Logger.Warn("failed to establish database connection, backing off...", "reason", err, "backoff_seconds", backoff.Seconds())
			time.Sleep(backoff)
			backoff = boFormula(backoff)
			try++
		}
	}
}

func (s *InventoryService) GetConn(ctx context.Context) *pgx.Conn {
	conninfo := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		s.Config.PgUser,
		s.Config.PgPass,
		s.Config.PgHost,
		s.Config.PgPort,
		s.Config.PgDbName,
		s.Config.PgSSLMode,
	)

	backoff := time.Second * 5 // secs
	time.Sleep(backoff)

	conn, err := tryConnect(s, ctx, conninfo, backoff, 5, func(backoff time.Duration) time.Duration {
		backoff = backoff + time.Second*5
		return backoff
	})

	if err != nil {
		return nil
	} else {
		s.Logger.Error(err)
	}

	return &conn
}

func (s *InventoryService) SellItem(ctx context.Context, req *pb.SellItemRequest) (*pb.SellItemResponse, error) {
	s.Logger.Warn("unimplemented method called.", "method", "SellItem")
	return &pb.SellItemResponse{Success: false, Message: "unimplemented"}, nil
}

func (s *InventoryService) GetBalance(ctx context.Context, req *pb.InventoryBalanceRequest) (*pb.InventoryBalanceResponse, error) {
	s.Logger.Warn("unimplemented method called.", "method", "SellItem")
	return &pb.InventoryBalanceResponse{Gold: 0, Tokens: 0}, nil
}

func (s *InventoryService) AddItemsToBackpack(ctx context.Context, req *pb.AddItemsToBackpackRequest) (*pb.AddItemsToBackpackResponse, error) {
	builder := sq.Insert("character_inventory_item_instance").
		Columns("character_id, item_instance_id").
		PlaceholderFormat(sq.Dollar)

	for _, itemInstanceId := range req.ItemInstanceIds {
		builder = builder.Values(req.GetCharacterId(), itemInstanceId)
	}

	sql, args, err := builder.ToSql()
	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	_, err = s.DbPool.Exec(ctx, sql, args...)
	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	return &pb.AddItemsToBackpackResponse{Success: true}, nil
}

func (s *InventoryService) ListBackpackItems(ctx context.Context, req *pb.ListBackpackItemsRequest) (*pb.ListBackpackItemsResponse, error) {
	cte, _, err := sq.Select("cte.*").
		From("character_inventory_item_instance AS cte").
		Where("cte.character_id = $1").
		ToSql()

	if err != nil {
		s.Logger.Error("failed to prepare cte", "err", err)
		return nil, err
	}

	prefix := fmt.Sprintf("WITH cte AS (%s)", cte)

	sql, _, err := sq.Select("templ.*, es.*, it.*, i.*").
		From("cte AS cte").
		InnerJoin("item_instances AS inst ON inst.item_instance_id = cte.item_instance_id").
		InnerJoin("items AS i ON i.id = inst.item_id").
		InnerJoin("item_templates AS templ ON templ.id = i.item_template_id").
		InnerJoin("item_types AS it ON it.id = templ.item_type_id").
		InnerJoin("equip_slots AS es ON es.id = templ.equip_slot_id").
		ToSql()

	if err != nil {
		s.Logger.Error("failed to prepare sql", "err", err)
		return nil, err
	}

	sql = fmt.Sprintf("%s %s", prefix, sql)
	rows, err := s.DbPool.Query(ctx, sql, req.CharacterId)
	if err != nil {
		s.Logger.Error("failed to run query", "err", err)
		return nil, err
	}

	var items []*pbRepo.Item
	for rows.Next() {
		var foo *any
		item := &pbRepo.Item{}
		itemTemplate := &pbRepo.ItemTemplate{}
		equipSlot := &pbRepo.EquipSlot{}
		itemType := &pbRepo.ItemType{}

		err := rows.Scan(
			&itemTemplate.Id,
			&itemTemplate.Name,
			&itemTemplate.ItemTypeId,
			&itemTemplate.EquipSlotId,
			&itemTemplate.Description,
			&itemTemplate.GoldPrice,
			&itemTemplate.BuyableWithTokens,
			&itemTemplate.TokenPrice,
			&foo, // metadata, unnecessary atm
			&equipSlot.Id,
			&equipSlot.Code,
			&equipSlot.Name,
			&itemType.Id,
			&itemType.Code,
			&itemType.Name,
			&item.Id,
			&item.Prefix,
			&item.Suffix,
			&item.ItemTemplateId,
			&item.Health,
			&item.Power,
			&item.Strength,
			&item.Spellpower,
			&item.BonusDamage,
			&item.BonusArmor,
		)

		if err != nil {
			s.Logger.Error(err)
			return nil, err
		}

		item.ItemTemplate = itemTemplate
		item.ItemTemplate.ItemType = itemType
		item.ItemTemplate.EquipSlot = equipSlot

		items = append(items, item)
	}

	return &pb.ListBackpackItemsResponse{Items: items}, nil
}

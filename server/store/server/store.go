package server

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	pbArmory "github.com/komadiina/spelltext/proto/armory"
	pbInventory "github.com/komadiina/spelltext/proto/inventory"
	pb "github.com/komadiina/spelltext/proto/store"
	"github.com/komadiina/spelltext/server/store/config"
	"github.com/komadiina/spelltext/utils/singleton/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type StoreService struct {
	pb.UnimplementedStoreServer
	DbPool *pgxpool.Pool
	Config *config.Config
	Logger *logging.Logger
}

func tryConnect(s *StoreService, context context.Context, conninfo string, backoff time.Duration, maxRetries int, boFormula func(time.Duration) time.Duration) (pgx.Conn, error) {
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

func (s *StoreService) GetConn(ctx context.Context) *pgx.Conn {
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
	}

	return &conn
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

	var vendors []*pb.Vendor
	for rows.Next() {
		v := &pb.Vendor{}
		err := rows.Scan(&v.VendorId, &v.VendorName, &v.VendorWareDescription)

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
	prefix := fmt.Sprintf("WITH v_filt AS (%s)", cteSql)

	query := sq.
		Select("templ.id, templ.name, templ.item_type_id," +
			"templ.description, templ.gold_price, templ.buyable_with_tokens," +
			"COALESCE(i.prefix, '') AS prefix," +
			"COALESCE(i.suffix, '') AS suffix," +
			"v_filt.code, COALESCE(templ.token_price, 0)," +
			"COALESCE(i.health, 0) AS health," +
			"COALESCE(i.power, 0) AS power," +
			"COALESCE(i.strength, 0) AS strength," +
			"COALESCE(i.spellpower, 0) AS spellpower," +
			"COALESCE(i.bonus_damage, 0) AS bonus_damage," +
			"COALESCE(i.bonus_armor, 0) AS bonus_armor").
		From("item_templates AS templ").
		InnerJoin("v_filt ON v_filt.item_type_id = templ.item_type_id").
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

	var items []*pb.Item
	for rows.Next() {
		it := &pb.Item{}
		err := rows.Scan(
			&it.Id,
			&it.Name,
			&it.ItemTypeId,
			&it.Description,
			&it.GoldPrice,
			&it.BuyableWithTokens,
			&it.Prefix,
			&it.Suffix,
			&it.ItemTypeCode,
			&it.TokenPrice,
			&it.Health,
			&it.Power,
			&it.Strength,
			&it.Spellpower,
			&it.BonusDamage,
			&it.BonusArmor,
		)

		if err != nil {
			s.Logger.Error("failed to scan", "err", err)
			return nil, err
		}

		items = append(items, it)
	}

	return &pb.ListVendorItemResponse{Items: items, TotalCount: -1}, nil
}

func (s *StoreService) AddItem(ctx context.Context, req *pb.AddItemRequest) (*pb.AddItemResponse, error) {
	s.Logger.Warn("unimplemented method called (AddItem)")
	return nil, nil
}

func (s *StoreService) BuyItems(ctx context.Context, req *pb.BuyItemRequest) (*pb.BuyItemResponse, error) {
	s.Logger.Debugf("BuyItems(%+v)", req)

	// tx, err := s.DbPool.Begin(ctx)
	// defer func() {_ = tx.Rollback(ctx)}()

	// if err != nil {
	// 	s.Logger.Error("failed to begin transaction", "err", err)
	// 	return nil, err
	// }
	
	errResp := &pb.BuyItemResponse{Success: false, Message: "error ocurred"}
	// check if character has enough gold
	sql, _, err := sq.Select("c.gold").
		From("characters AS c").
		Where("c.character_id = $1").
		ToSql()

	c := &pbArmory.TCharacter{Id: req.CharacterId}
	rows, err := s.DbPool.Query(ctx, sql, req.CharacterId)
	// rows, err := tx.Query(ctx, sql, req.CharacterId)
	for rows.Next() {
		err := rows.Scan(&c.Gold)

		if err != nil {
			s.Logger.Error(err)
			return errResp, err
		}
	}

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
	// rows, err = tx.Query(ctx, sql, args...)
	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}
	
	var items []*pb.Item
	for rows.Next() {
		item := &pb.Item{}
		
		err := rows.Scan(&item.Id, &item.GoldPrice, &item.ItemTemplateId) // repeated but skip block of code
		if err != nil {
			s.Logger.Error(err)
			return errResp, err
		}

		items = append(items, item)
	}

	var sum uint64 = 0
	for _, item := range items {
		sum += item.GetGoldPrice()
	}

	if sum > c.Gold {
		s.Logger.Infof("character %s overbought attempt, gold_amount=%d, character_gold=%d", c.Name, sum, c.Gold)
		return &pb.BuyItemResponse{Success: false, Message: "not enough gold"}, fmt.Errorf("overbuy attempt")
	}

	var itemInstanceIds []uint64
	sql = "INSERT INTO item_instances VALUES (DEFAULT, $1, $2, DEFAULT, DEFAULT) RETURNING item_instance_id" // could? refactor into single multi-line insert
	for _, item := range items {
		rows, err := s.DbPool.Query(ctx, sql, item.GetId(), c.Id)
		// rows, err := tx.Query(ctx, sql, item.GetId(), c.Id)
		if err != nil {
			s.Logger.Error(err)
			return errResp, err
		}

		for rows.Next() {
			var id uint64
			err := rows.Scan(&id)
			if err != nil {
				s.Logger.Error(err)
				return errResp, err
			}

			itemInstanceIds = append(itemInstanceIds, id)
		}
	}

	s.Logger.Info(itemInstanceIds)
	
	// update character gold
	sql = "UPDATE characters SET gold = gold - $1 WHERE character_id = $2"
	_, err = s.DbPool.Exec(ctx, sql, sum, c.Id)
	// _, err = tx.Exec(ctx, sql, sum, c.Id)
	if err != nil {
		s.Logger.Error(err)
		return errResp, err
	}

	// TODO: move from direct service-service to MQ (problem: wait-for-ack)
	conn, err := grpc.NewClient("inventoryserver:50053", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		s.Logger.Error("failed to init inventory client", "reason", err)
		return errResp, err
	}
	invClient := pbInventory.NewInventoryClient(conn)
	_, err = invClient.AddItemsToBackpack(ctx, &pbInventory.AddItemsToBackpackRequest{
		CharacterId: c.Id,
		ItemInstanceIds: itemInstanceIds, // ERROR: using ItemIds instead of last inserted ids into 'item_instances' table
	})

	if err != nil {
		s.Logger.Error(err)
		return errResp, err
	}

	// err = tx.Commit(ctx)
	// if err != nil {
	// 	s.Logger.Error(err)
	// 	return errResp, err
	// }

	return &pb.BuyItemResponse{Success: true, Message: "items bought & added to inventory"}, nil
}

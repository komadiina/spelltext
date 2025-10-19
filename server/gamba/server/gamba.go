package server

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	pb "github.com/komadiina/spelltext/proto/gamba"
	pbInventory "github.com/komadiina/spelltext/proto/inventory"
	pbRepo "github.com/komadiina/spelltext/proto/repo"
	"github.com/komadiina/spelltext/server/gamba/config"
	"github.com/komadiina/spelltext/utils/singleton/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GambaService struct {
	pb.UnimplementedGambaServer
	DbPool *pgxpool.Pool
	Config *config.Config
	Logger *logging.Logger
}

func tryConnect(s *GambaService, context context.Context, conninfo string, backoff time.Duration, maxRetries int, boFormula func(time.Duration) time.Duration) (pgx.Conn, error) {
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

func (s *GambaService) GetConn(ctx context.Context) *pgx.Conn {
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

func (s *GambaService) GetChests(ctx context.Context, req *pb.GetChestsRequest) (*pb.GetChestsResponse, error) {
	sql, _, err := sq.
		Select("*").
		From("gamba_chests").
		ToSql()

	if err != nil {
		s.Logger.Error("failed to build sql", "err", err)
		return nil, err
	}

	rows, err := s.DbPool.Query(ctx, sql)
	if err != nil {
		s.Logger.Error("failed to run query", "err", err)
		return nil, err
	}

	var gchests []*pbRepo.GambaChest
	for rows.Next() {
		gc := &pbRepo.GambaChest{}

		err := rows.Scan(
			&gc.Id,
			&gc.Name,
			&gc.Description,
			&gc.GoldPrice,
		)

		if err != nil {
			s.Logger.Error(err)
			return nil, err
		}

		gchests = append(gchests, gc)
	}

	return &pb.GetChestsResponse{Chests: gchests}, nil
}

func (s *GambaService) OpenChest(ctx context.Context, req *pb.OpenChestRequest) (*pb.OpenChestResponse, error) {
	// check if character has enough gold
	prefix := "WITH gc_filt AS (SELECT price FROM gamba_chests WHERE id = $1)"
	sql, _, err := sq.
		Select("c.gold, (SELECT price FROM gamba_chests WHERE id = $1) AS price").
		From("characters AS c").
		Where("c.character_id = $2").
		Limit(1).
		ToSql()

	if err != nil {
		s.Logger.Error("failed to build query", "error", err)
		return nil, err
	}

	sql = prefix + " " + sql

	rows, err := s.DbPool.Query(ctx, sql, req.ChestId, req.CharacterId)
	if err != nil {
		s.Logger.Error("failed to run query", "err", err)
		return nil, err
	}

	char := &pbRepo.Character{}
	gc := &pbRepo.GambaChest{}
	for rows.Next() {
		err := rows.Scan(&char.Gold, &gc.GoldPrice)
		if err != nil {
			s.Logger.Error(err)
			return nil, err
		}
	}

	if char.Gold < gc.GoldPrice {
		err = fmt.Errorf("not enough gold, have=%d, need=%d", char.Gold, gc.GoldPrice)
		return nil, err
	}

	// get list of items by chest.id
	sql, _, err = sq.
		Select("i.*, it.*").
		From("gamba_chest_contents AS gcc").
		InnerJoin("items AS i ON i.id = gcc.item_id").
		InnerJoin("item_templates AS it ON it.id = i.item_template_id").
		Where("gcc.gamba_chest_id = $1").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		s.Logger.Error("failed to build query", "reason", err)
		return nil, err
	}

	var items []*pbRepo.Item
	rows, err = s.DbPool.Query(ctx, sql, req.ChestId)
	for rows.Next() {
		var foo *any
		it := &pbRepo.ItemTemplate{}
		i := &pbRepo.Item{}

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
			&it.Id,
			&it.Name,
			&it.ItemTypeId,
			&it.EquipSlotId,
			&it.Description,
			&it.GoldPrice,
			&it.BuyableWithTokens,
			&it.TokenPrice,
			&foo,
		)

		if err != nil {
			s.Logger.Error(err)
			return nil, err
		}

		i.ItemTemplate = it

		items = append(items, i)
	}

	reward := items[rand.Intn(len(items)-1)]

	// update character gold
	sql = "UPDATE characters SET gold = gold - $1 WHERE character_id = $2"
	_, err = s.DbPool.Exec(ctx, sql, gc.GoldPrice, req.CharacterId)
	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	// create a new item instance
	sql = "INSERT INTO item_instances (item_id, owner_character_id) VALUES ($1, $2) RETURNING item_instance_id"
	rows, err = s.DbPool.Query(ctx, sql, reward.GetId(), req.GetCharacterId())
	var instanceId uint64 = 0
	for rows.Next() {
		err = rows.Scan(&instanceId)
		if err != nil {
			s.Logger.Error(err)
			return nil, err
		}
	}

	// contact inventoryserver
	conn, err := grpc.NewClient("inventoryserver:50053", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		s.Logger.Error("failed to init inventory client", "reason", err)
		return nil, err
	}

	client := pbInventory.NewInventoryClient(conn)
	_, err = client.AddItemsToBackpack(ctx,
		&pbInventory.AddItemsToBackpackRequest{
			CharacterId:     req.CharacterId,
			ItemInstanceIds: []uint64{instanceId},
		},
	)

	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	return &pb.OpenChestResponse{Item: reward}, nil
}

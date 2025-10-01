package server

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	pb "github.com/komadiina/spelltext/proto/store"
	"github.com/komadiina/spelltext/server/store/config"
	"github.com/komadiina/spelltext/utils/singleton/logging"
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

func (s *StoreService) ListItems(ctx context.Context, req *pb.StoreListItemRequest) (*pb.ItemListResponse, error) {
	sql := "SELECT * FROM item_templates"
	rows, err := s.DbPool.Query(ctx, sql)
	if err != nil {
		logging.Get("store").Error("failed to query", "reason", err)
		return nil, err
	}

	var items []*pb.Item
	for rows.Next() {
		it := &pb.Item{}
		err := rows.Scan(&it.Id, &it.Name, &it.ItemTypeId, &it.Rarity, &it.Stackable, &it.StackSize, &it.BindOnPickup, &it.Description, &it.Metadata)
		if err != nil {
			logging.Get("store").Error("failed to scan", "reason", err)
			return nil, err
		}
		items = append(items, it)
	}

	return &pb.ItemListResponse{Items: items, TotalCount: 0}, nil
}

func (s *StoreService) AddItem(ctx context.Context, req *pb.AddItemRequest) (*pb.AddItemResponse, error) {
	panic("unimplemented")
}

func (s *StoreService) BuyItem(ctx context.Context, req *pb.BuyItemRequest) (*pb.BuyItemResponse, error) {
	panic("unimplemented")
}

package server

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	pb "github.com/komadiina/spelltext/proto/inventory"
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

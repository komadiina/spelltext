package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	pb "github.com/komadiina/spelltext/proto/store"
	"github.com/komadiina/spelltext/server/store/config"
	"github.com/komadiina/spelltext/server/store/server"
	"github.com/komadiina/spelltext/utils/singleton/logging"
	"google.golang.org/grpc"
)

func InitializePool(s *server.StoreService, context context.Context, conninfo string, backoff time.Duration, maxRetries int, boFormula func(time.Duration) time.Duration) error {
	try := 1
	for {
		conn, err := pgx.Connect(context, conninfo)

		if err != nil && try >= maxRetries {
			// conn not established, max retries exceeded
			s.Logger.Fatal(err)
		} else if err == nil && try < maxRetries {
			// conn established within maxRetries
			s.Logger.Info("pgpool connection established, creating pool..")
			conn.Close(context)

			pool, err := pgxpool.New(context, fmt.Sprintf(
				"postgres://%s:%s@%s:%d/%s?sslmode=%s",
				s.Config.PgUser,
				s.Config.PgPass,
				s.Config.PgHost,
				s.Config.PgPort,
				s.Config.PgDbName,
				s.Config.PgSSLMode,
			))

			if err != nil {
				log.Fatal("unable to create pool", "reason", err)
			} else {
				log.Info("pgxpool (dpool, via pgpool-ii) initialized")
			}

			s.DbPool = pool

			return nil
		} else if err != nil && try < maxRetries {
			// conn not established, backoff
			s.Logger.Warn("failed to establish database connection, backing off...", "reason", err, "backoff_seconds", backoff.Seconds())
			time.Sleep(backoff)
			backoff = boFormula(backoff)
			try++
		}
	}
}

func main() {
	ctx := context.Background()

	logging.Init(log.InfoLevel, "storeserver", false)
	logger := logging.Get("storeserver", false)
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal(err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.ServicePort))
	if err != nil {
		logger.Error("failed to listen", err)
		os.Exit(1)
	}

	s := grpc.NewServer()
	ss := server.StoreService{Config: cfg, Logger: logger}

	conninfo := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.PgUser,
		cfg.PgPass,
		cfg.PgHost,
		cfg.PgPort,
		cfg.PgDbName,
		cfg.PgSSLMode,
	)
	err = InitializePool(&ss, ctx, conninfo, time.Second*5, 10, func(bo time.Duration) time.Duration {
		return bo + time.Second*5
	})

	if err != nil {
		ss.Logger.Fatal("failed to connect to database/initialize pgxpool, not serving.", "reason", err)
	}

	pb.RegisterStoreServer(s, &ss)
	logger.Info(fmt.Sprintf("%s v%s listening on %s:%d", "storeserver", "0.3.0", "127.0.0.1", ss.Config.ServicePort))

	if err := s.Serve(lis); err != nil {
		logger.Error("failed to serve", "reason", err)
		os.Exit(1)
	}
}

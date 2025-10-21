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
	pb "github.com/komadiina/spelltext/proto/gamba"
	"github.com/komadiina/spelltext/server/gamba/config"
	"github.com/komadiina/spelltext/server/gamba/server"
	"github.com/komadiina/spelltext/utils/singleton/logging"
	"google.golang.org/grpc"
)

const banner = `
                _ _ _            _   
               | | | |          | |  
 ___ _ __   ___| | | |_ _____  _| |_ 
/ __| '_ \ / _ \ | | __/ _ \ \/ / __|
\__ \ |_) |  __/ | | ||  __/>  <| |_ 
|___/ .__/ \___|_|_|\__\___/_/\_\\__|
    | |                              
    |_|                              

`

func InitializePool(s *server.GambaService, context context.Context, conninfo string, backoff time.Duration, maxRetries int, boFormula func(time.Duration) time.Duration) error {
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
				"user=%s password=%s host=%s port=%d dbname=%s sslmode=%s",
				s.Config.PgUser,
				s.Config.PgPass,
				s.Config.PgHost,
				s.Config.PgPort,
				s.Config.PgDbName,
				s.Config.PgSSLMode,
			))

			if err != nil {
				s.Logger.Fatal("unable to create pool", "reason", err)
			} else {
				s.Logger.Info("pgxpool (dpool, via pgpool-ii) initialized")
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

var version = os.Getenv("VERSION")

func main() {
	ctx := context.Background()

	logging.Init(log.InfoLevel, "gambaserver", false)
	logger := logging.Get("gambaserver", false)

	logger.Infof(`%s%sversion=%s`, banner, "\n", version)

	logger.Info("loading config...", "CONFIG_FILE", os.Getenv("CONFIG_FILE"))
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal(err)
	} else {
		logger.Info("gambaserver config loaded.")
		logger.Infof("conninfo=%v:%v@%v:%v/%v?sslMode=%v, port=%d",
			cfg.PgUser, cfg.PgPass, cfg.PgHost, cfg.PgPort, cfg.PgDbName, cfg.PgSSLMode, cfg.ServicePort)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.ServicePort))
	if err != nil {
		logger.Error("failed to listen", err)
		os.Exit(1)
	}

	s := grpc.NewServer()
	ss := server.GambaService{Config: cfg, Logger: logger}

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

	defer func() {
		ss.Logger.Info("closing pgx dbconn pool...")
		ss.DbPool.Close()
	}()

	if err != nil {
		ss.Logger.Fatal("failed to connect to database/initialize pgxpool, not serving.", "reason", err)
	}

	pb.RegisterGambaServer(s, &ss)
	logger.Info(fmt.Sprintf("%s v%s listening on %s:%d", "gambaserver", "0.3.0", "127.0.0.1", ss.Config.ServicePort))

	if err := s.Serve(lis); err != nil {
		logger.Error("failed to serve", "reason", err)
		os.Exit(1)
	}
}

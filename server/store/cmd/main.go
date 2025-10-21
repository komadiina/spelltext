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
	pbInventory "github.com/komadiina/spelltext/proto/inventory"
	pb "github.com/komadiina/spelltext/proto/store"
	"github.com/komadiina/spelltext/server/store/config"
	"github.com/komadiina/spelltext/server/store/server"
	"github.com/komadiina/spelltext/utils/singleton/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

func InitClientConn(target string, credentials grpc.DialOption, backoff int, maxRetries int) (*grpc.ClientConn, error) {
	try := 1
	for {
		conn, err := grpc.NewClient(target, credentials)

		if err != nil && try >= maxRetries {
			return nil, err
		} else if err == nil && try < maxRetries {
			return conn, nil
		} else if err != nil && try < maxRetries {
			backoff *= 3
			time.Sleep(time.Duration(backoff) * time.Second)
			try++
		}
	}
}

var version = os.Getenv("VERSION")

func main() {
	ctx := context.Background()

	logging.Init(log.InfoLevel, "storeserver", false)
	logger := logging.Get("storeserver", false)

	logger.Infof("%s\n%s", banner, version)

	logger.Info("loading config...", "CONFIG_FILE", os.Getenv("CONFIG_FILE"))
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal(err)
	} else {
		logger.Info("storeserver config loaded.")
		logger.Infof("conninfo=%v:%v@%v:%v/%v?sslMode=%v, port=%d",
			cfg.PgUser, cfg.PgPass, cfg.PgHost, cfg.PgPort, cfg.PgDbName, cfg.PgSSLMode, cfg.ServicePort)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.ServicePort))
	if err != nil {
		logger.Error("failed to listen", err)
		os.Exit(1)
	}

	s := grpc.NewServer()
	ss := server.StoreService{Config: cfg, Logger: logger}

	// init inventory clientConn
	clientConn, err := InitClientConn("inventoryserver:50053", grpc.WithTransportCredentials(insecure.NewCredentials()), 5, 10)
	if err != nil {
		logger.Fatal(err)
	}

	ss.Connections = &server.Connections{
		Inventory: clientConn,
	}
	ss.Clients = &server.Clients{
		Inventory: pbInventory.NewInventoryClient(clientConn),
	}
	defer ss.Connections.Inventory.Close()

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

	pb.RegisterStoreServer(s, &ss)
	logger.Info(fmt.Sprintf("%s v%s listening on %s:%d", "storeserver", "0.3.0", "127.0.0.1", ss.Config.ServicePort))

	if err := s.Serve(lis); err != nil {
		logger.Error("failed to serve", "reason", err)
		os.Exit(1)
	}
}

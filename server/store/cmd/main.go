package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/charmbracelet/log"
	pbInventory "github.com/komadiina/spelltext/proto/inventory"
	pb "github.com/komadiina/spelltext/proto/store"
	"github.com/komadiina/spelltext/server/store/config"
	"github.com/komadiina/spelltext/server/store/db"
	"github.com/komadiina/spelltext/server/store/health"
	"github.com/komadiina/spelltext/server/store/server"
	"github.com/komadiina/spelltext/server/store/services"
	"github.com/komadiina/spelltext/shared"
	"github.com/komadiina/spelltext/utils/singleton/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

var version = os.Getenv("VERSION")

func main() {
	ctx := context.Background()

	logging.Init(log.InfoLevel, "storeserver", false)
	logger := logging.Get("storeserver", false)

	logger.Infof("%s\n%s", shared.BANNER, version)

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
	target := fmt.Sprint("inventoryserver:", cfg.InventoryServicePort)
	clientConn, err := services.InitClientConn(
		logger,
		target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		5,
		10)
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

	go health.InitMonitor(
		&ss,
		target,
		ss.Clients.Inventory,
		func(ss *server.StoreService, cc *grpc.ClientConn) {
			ss.Connections.Inventory = cc
			ss.Logger.Infof("server is back up, healthy. service=%s", target)
		}).
		Run(ctx)

	conninfo := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.PgUser,
		cfg.PgPass,
		cfg.PgHost,
		cfg.PgPort,
		cfg.PgDbName,
		cfg.PgSSLMode,
	)
	err = db.InitializePool(&ss, ctx, conninfo, time.Second*5, 10, func(bo time.Duration) time.Duration {
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

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		logger.Error("failed to serve", "reason", err)
		os.Exit(1)
	}
}

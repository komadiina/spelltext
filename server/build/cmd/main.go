package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/charmbracelet/log"
	pb "github.com/komadiina/spelltext/proto/build"
	"github.com/komadiina/spelltext/proto/char"
	"github.com/komadiina/spelltext/server/build/config"
	"github.com/komadiina/spelltext/server/build/db"
	"github.com/komadiina/spelltext/server/build/health"
	"github.com/komadiina/spelltext/server/build/server"
	"github.com/komadiina/spelltext/server/build/services"
	"github.com/komadiina/spelltext/shared"
	"github.com/komadiina/spelltext/utils/singleton/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var version = os.Getenv("VERSION")

func main() {
	ctx := context.Background()

	logging.Init(log.InfoLevel, "buildserver", false)
	logger := logging.Get("buildserver", false)

	logger.Infof("%s\n%s", shared.BANNER, version)

	logger.Info("loading config...", "CONFIG_FILE", os.Getenv("CONFIG_FILE"))
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal(err)
	} else {
		logger.Info("buildserver config loaded.")
		logger.Infof("conninfo=%v:%v@%v:%v/%v?sslMode=%v, port=%d",
			cfg.PgUser, cfg.PgPass, cfg.PgHost, cfg.PgPort, cfg.PgDbName, cfg.PgSSLMode, cfg.ServicePort)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.ServicePort))
	if err != nil {
		logger.Error("failed to listen", err)
		os.Exit(1)
	}

	s := grpc.NewServer()
	ss := server.BuildService{Config: cfg, Logger: logger}

	target := fmt.Sprint("charserver:", cfg.CharacterServicePort)
	conn, err := services.InitClientConn(
		logger,
		target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		5,
		10,
	)
	ss.Clients = &server.Clients{
		Character: char.NewCharacterClient(conn),
	}
	ss.Connections = &server.Connections{
		Character: conn,
	}
	defer conn.Close()

	go health.InitMonitor(
		&ss,
		target,
		ss.Clients.Character,
		func(bs *server.BuildService, cc *grpc.ClientConn) {
			bs.Clients.Character = char.NewCharacterClient(cc)
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

	pb.RegisterBuildServer(s, &ss)
	logger.Info(fmt.Sprintf("%s v%s listening on %s:%d", "buildserver", "0.5.0", "127.0.0.1", ss.Config.ServicePort))

	if err := s.Serve(lis); err != nil {
		logger.Error("failed to serve", "reason", err)
		os.Exit(1)
	}
}

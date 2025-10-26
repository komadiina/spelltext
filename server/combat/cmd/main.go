package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/charmbracelet/log"
	pbChar "github.com/komadiina/spelltext/proto/char"
	pb "github.com/komadiina/spelltext/proto/combat"
	"github.com/komadiina/spelltext/server/combat/config"
	"github.com/komadiina/spelltext/server/combat/db"
	"github.com/komadiina/spelltext/server/combat/health"
	"github.com/komadiina/spelltext/server/combat/server"
	"github.com/komadiina/spelltext/server/combat/services"
	"github.com/komadiina/spelltext/shared"
	"github.com/komadiina/spelltext/utils/singleton/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var version = os.Getenv("VERSION")

func main() {
	ctx := context.Background()

	logging.Init(log.InfoLevel, "combatserver", false)
	logger := logging.Get("combatserver", false)

	logger.Infof(`%s%sversion=%s`, shared.BANNER, "\n", version)

	logger.Info("loading config...", "CONFIG_FILE", os.Getenv("CONFIG_FILE"))
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal(err)
	} else {
		logger.Info("combatserver config loaded.")
		logger.Infof("conninfo=%v:%v@%v:%v/%v?sslMode=%v, port=%d",
			cfg.PgUser, cfg.PgPass, cfg.PgHost, cfg.PgPort, cfg.PgDbName, cfg.PgSSLMode, cfg.ServicePort)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.ServicePort))
	if err != nil {
		logger.Error("failed to listen", err)
		os.Exit(1)
	}

	s := grpc.NewServer()
	ss := server.CombatService{Config: cfg, Logger: logger}

	conninfo := fmt.Sprintf(
		"user=%s password=%s host=%s port=%d dbname=%s sslmode=%s",
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

	target := fmt.Sprintf("charserver:%d", cfg.CharacterServicePort)
	clientConn, err := services.InitClientConn(
		logger,
		target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		5,
		10,
	)
	if err != nil {
		logger.Fatal(err)
	}

	ss.Connections = &server.Connections{
		Character: clientConn,
	}

	ss.Clients = &server.Clients{
		Character: pbChar.NewCharacterClient(clientConn),
	}

	defer ss.Connections.Character.Close()

	go health.InitMonitor(
		&ss,
		target,
		ss.Clients.Character,
		func(s *server.CombatService, conn *grpc.ClientConn) {
			s.Clients.Character = pbChar.NewCharacterClient(conn)
			ss.Logger.Infof("service is back up, healthy. service=%s", target)
		}).
		Run(ctx)

	defer func() {
		ss.Logger.Info("closing pgx dbconn pool...")
		ss.DbPool.Close()
	}()

	if err != nil {
		ss.Logger.Fatal("failed to connect to database/initialize pgxpool, not serving.", "reason", err)
	}

	pb.RegisterCombatServer(s, &ss)
	logger.Info(fmt.Sprintf("%s v%s listening on %s:%d", "combatserver", version, "127.0.0.1", ss.Config.ServicePort))

	if err := s.Serve(lis); err != nil {
		logger.Error("failed to serve", "reason", err)
		os.Exit(1)
	}
}

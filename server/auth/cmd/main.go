package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/charmbracelet/log"
	pb "github.com/komadiina/spelltext/proto/auth"
	pbChar "github.com/komadiina/spelltext/proto/char"
	"github.com/komadiina/spelltext/server/auth/config"
	"github.com/komadiina/spelltext/server/auth/db"
	"github.com/komadiina/spelltext/server/auth/health"
	"github.com/komadiina/spelltext/server/auth/server"
	"github.com/komadiina/spelltext/server/auth/services"
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

var version = os.Getenv("VERSION")

func main() {
	ctx := context.Background()

	logging.Init(log.InfoLevel, "authserver", false)
	logger := logging.Get("authserver", false)

	logger.Infof("%s\n%s", banner, version)

	logger.Info("loading config...", "CONFIG_FILE", os.Getenv("CONFIG_FILE"))
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal(err)
	} else {
		logger.Infof("authserver config loaded: %+v", cfg)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.ServicePort))
	if err != nil {
		logger.Error("failed to listen", err)
		os.Exit(1)
	}

	s := grpc.NewServer()
	ss := server.AuthService{Config: cfg, Logger: logger}

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
		func(s *server.AuthService, conn *grpc.ClientConn) {
			s.Clients.Character = pbChar.NewCharacterClient(conn)
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

	pb.RegisterAuthServer(s, &ss)
	logger.Info(fmt.Sprintf("%s v%s listening on %s:%d", "authserver", version, "127.0.0.1", ss.Config.ServicePort))

	if err := s.Serve(lis); err != nil {
		logger.Error("failed to serve", "reason", err)
		os.Exit(1)
	}
}

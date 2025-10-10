package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/charmbracelet/log"
	pb "github.com/komadiina/spelltext/proto/chat"
	"github.com/komadiina/spelltext/server/chat/config"
	"github.com/komadiina/spelltext/server/chat/server"
	"github.com/komadiina/spelltext/utils/singleton/logging"
	"google.golang.org/grpc"
)

var version = os.Getenv("VERSION")

func main() {
	logging.Init(log.InfoLevel, "chatserver", false)
	logger := logging.Get("chatserver", false)

	logger.Infof(`
		// -------------------- //
		// ---- chatserver ---- //
		// ----   %v   ---- //
		// -------------------- //`, version)

	logger.Info("loading config...", "CONFIG_FILE", os.Getenv("CONFIG_FILE"))
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal(err)
	} else {
		logger.Info("chatserver config loaded.")
		logger.Infof("nats_url=%v, port=%v, max_async_publish=%v",
			cfg.NatsURL, cfg.Port, cfg.MaxAsyncPublish)
	}

	logger.Info("initializing nats..")
	nc, err := server.InitNats()
	defer nc.Drain()

	if err != nil {
		log.Fatal(err)
	}

	logger.Info(fmt.Sprintf("Starting %s v%s...", "chatserver", version))
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		logger.Error("failed to listen", err)
		os.Exit(1)
	}

	s := grpc.NewServer()
	pb.RegisterChatServer(s, &server.ChatService{Nats: nc, Config: cfg, Logger: logger})

	logger.Infof("chatserver v%s listening on localhost:%d", version, cfg.Port)

	if err := s.Serve(lis); err != nil {
		logger.Error("failed to serve", "reason", err)
		os.Exit(1)
	}
}

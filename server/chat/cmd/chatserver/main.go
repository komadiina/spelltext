package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/charmbracelet/log"
	pb "github.com/komadiina/spelltext/proto/chat"
	"github.com/komadiina/spelltext/server/chat/server"
	"github.com/komadiina/spelltext/utils/singleton/logging"
	"google.golang.org/grpc"
)

const (
	version = "0.0.1"
)

var (
	port = flag.Int("port", 50051, "port to listen on")
	name = flag.String("name", "chatserver", "server name")
	addr = flag.String("addr", "localhost", "server address")
)

func main() {
	fmt.Println("Starting server...")

	logging.Init(log.InfoLevel)
	logger := logging.Get()

	logger.Info(fmt.Sprintf("Starting %s v%s...", *name, version))

	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		logger.Error("failed to listen", err)
		os.Exit(1)
	}

	s := grpc.NewServer()
	pb.RegisterChatServiceServer(s, &server.ChatService{})

	logger.Info(fmt.Sprintf("%s v%s listening on %s:%d", *name, version, *addr, *port))

	if err := s.Serve(lis); err != nil {
		logger.Error("failed to serve", err)
		os.Exit(1)
	}
}

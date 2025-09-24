package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/charmbracelet/log"
	pb "github.com/komadiina/spelltext/proto/chat"
	"google.golang.org/grpc"
)

const (
	version = "0.0.1"
)

var (
	port = flag.Int("port", 50051, "port to listen on")
	name = flag.String("name", "chatserver", "server name")
	addr = flag.String("addr", "127.0.0.1", "server address")
)

var logger log.Logger

type ChatService struct {
	pb.UnimplementedChatServiceServer
	sentMessages int
}

func (s *ChatService) SendChatMessage(ctx context.Context, req *pb.SendChatMessageRequest) (*pb.SendChatMessageResponse, error) {
	logger.Info(fmt.Sprintf(fmt.Sprintf("Received message from %s: %s", req.Sender, req.Message)))

	s.sentMessages++
	return &pb.SendChatMessageResponse{
		Sender:  *name,
		Success: true,
	}, nil
}

func main() {
	logger = *log.NewWithOptions(os.Stdout, log.Options{
		ReportCaller:    true,
		ReportTimestamp: true,
		TimeFormat:      time.DateTime,
	})

	logger.Info(fmt.Sprintf("Starting %s v%s...", *name, version))

	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *addr, *port))
	if err != nil {
		logger.Error("failed to listen", err)
		os.Exit(1)
	}

	server := grpc.NewServer()
	pb.RegisterChatServiceServer(server, &ChatService{})

	logger.Info(fmt.Sprintf("%s v%s listening on %s:%d", *name, version, *addr, *port))

	if err := server.Serve(lis); err != nil {
		logger.Error("failed to serve", err)
		os.Exit(1)
	}
}

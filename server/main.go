package server

import (
	"context"
	"log"
	"net"

	pb "github.com/komadiina/spelltext/proto"
	"google.golang.org/grpc"
)

func (s *pb.Server) ServerPing(ctx context.Context, req *pb.Ping) (*pb.Pong, error) {
	return &pb.Pong{Message: "Hello " + req.Name}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	server := grpc.NewServer()

	err = server.Serve(lis)
	if err != nil {
		log.Fatal(err)
	}
}

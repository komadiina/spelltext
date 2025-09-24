package server

import (
	"context"
	"fmt"
	"log"

	pb "github.com/komadiina/spelltext/proto/chat"
	"github.com/komadiina/spelltext/utils/singleton/logging"
)

var logger log.Logger

type ChatService struct {
	pb.UnimplementedChatServiceServer
	sentMessages int
}

func (s *ChatService) SendChatMessage(ctx context.Context, req *pb.SendChatMessageRequest) (*pb.SendChatMessageResponse, error) {
	logger := logging.Get()
	logger.Info(fmt.Sprintf(fmt.Sprintf("[#%d] Received message from %s: %s", s.sentMessages, req.Sender, req.Message)))

	s.sentMessages++
	return &pb.SendChatMessageResponse{
		Sender:  "chatserver",
		Success: true,
	}, nil
}

// Should broadcast to all users
func (s *ChatService) JoinChatroomMessage(ctx context.Context, req *pb.JoinChatroomMessageRequest) (*pb.JoinChatroomMessageResponse, error) {
	return &pb.JoinChatroomMessageResponse{
		Success: true,
	}, nil
}

// Should broadcast to all users
func (s *ChatService) LeaveChatroomMessage(ctx context.Context, req *pb.LeaveChatroomMessageRequest) (*pb.LeaveChatroomMessageResponse, error) {
	return &pb.LeaveChatroomMessageResponse{
		Success: true,
	}, nil
}

package server

import (
	"context"
	"fmt"
	"io"
	"log"

	pb "github.com/komadiina/spelltext/proto/chat"
	"github.com/komadiina/spelltext/utils/singleton/logging"
)

var logger log.Logger

type ChatService struct {
	pb.UnimplementedChatServiceServer
	sentMessages int
	Hub          *Hub
}

// Should broadcast to all users
func (s *ChatService) JoinChatroom(ctx context.Context, req *pb.JoinChatroomMessageRequest) (*pb.JoinChatroomMessageResponse, error) {
	s.Hub.Broadcast(req.GetUsername(), fmt.Sprintf("%s joined the chatroom.", req.GetUsername()))
	logging.Get().Info(fmt.Sprintf("%s joined the chatroom.", req.GetUsername()))

	return &pb.JoinChatroomMessageResponse{
		Success: true,
	}, nil
}

// Should broadcast to all users
func (s *ChatService) LeaveChatroom(ctx context.Context, req *pb.LeaveChatroomMessageRequest) (*pb.LeaveChatroomMessageResponse, error) {
	s.Hub.Broadcast(req.GetUsername(), fmt.Sprintf("%s left the chatroom.", req.GetUsername()))

	return &pb.LeaveChatroomMessageResponse{
		Success: true,
	}, nil
}

func (s *ChatService) SendChatMessage(ctx context.Context, req *pb.SendChatMessageRequest) (*pb.SendChatMessageResponse, error) {
	logger := logging.Get()
	logger.Info(fmt.Sprintf(fmt.Sprintf("[#%d] Received message from %s: %s", s.sentMessages, req.Sender, req.Message)))

	msg := s.Hub.Broadcast(req.Sender, req.Message)
	return &pb.SendChatMessageResponse{Sender: msg.Sender, Timestamp: msg.Timestamp, Success: true}, nil
}

func (s *ChatService) Subscribe(req *pb.SubscribeRequest, stream pb.ChatService_SubscribeServer) error {
	id, ch := s.Hub.Add()
	defer s.Hub.Remove(id)

	for {
		select {
		case <-stream.Context().Done():
			return nil

		case m, ok := <-ch:
			if !ok {
				return nil
			}
			if err := stream.Send(m); err != nil {
				if err == io.EOF {
					return nil
				}
				return err
			}
		}
	}
}

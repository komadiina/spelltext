package server

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	pb "github.com/komadiina/spelltext/proto/chat"
	"github.com/komadiina/spelltext/server/chat/config"
	"github.com/komadiina/spelltext/utils/singleton/logging"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

type ChatService struct {
	pb.UnimplementedChatServer
	sentMessages int
	Nats         *nats.Conn
	Config       *config.Config
	Logger       *logging.Logger
}

func publishMessage(s *ChatService, msg *pb.ChatMessage) error {
	js, err := s.Nats.JetStream()
	if err != nil {
		return err
	}

	msg.Message = fmt.Sprintf("[%s]: %s\n", msg.Sender, msg.Message)

	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	m := &nats.Msg{
		Subject: "chat.global",
		Data:    data,
		Header:  nats.Header{"Content-Type": []string{"application/protobuf"}, "Nats-Msg-Uuid": []string{uuid.NewString()}},
	}

	_, err = js.PublishMsg(m)
	if err != nil {
		return err
	}

	return nil
}

// Should broadcast to all users
func (s *ChatService) JoinChatroom(ctx context.Context, req *pb.JoinChatroomMessageRequest) (*pb.JoinChatroomMessageResponse, error) {
	s.Logger.Info(fmt.Sprintf("%s joined the chatroom.", req.GetUsername()))
	err := publishMessage(s, &pb.ChatMessage{Sender: "chatserver", Message: fmt.Sprintf("%s joined the chatroom.", req.GetUsername())})

	if err != nil {
		return &pb.JoinChatroomMessageResponse{
			Success: false,
		}, err
	}

	return &pb.JoinChatroomMessageResponse{
		Success: true,
	}, nil
}

// Should broadcast to all users
func (s *ChatService) LeaveChatroom(ctx context.Context, req *pb.LeaveChatroomMessageRequest) (*pb.LeaveChatroomMessageResponse, error) {
	err := publishMessage(s, &pb.ChatMessage{Sender: req.GetUsername(), Message: fmt.Sprintf("%s left the chatroom.", req.GetUsername())})
	if err != nil {
		return &pb.LeaveChatroomMessageResponse{
			Success: false,
		}, err
	}

	return &pb.LeaveChatroomMessageResponse{
		Success: true,
	}, nil
}

func (s *ChatService) SendChatMessage(ctx context.Context, req *pb.SendChatMessageRequest) (*pb.SendChatMessageResponse, error) {
	s.Logger.Info(fmt.Sprintf(fmt.Sprintf("[#%d] Received message from %s: %s", s.sentMessages, req.Sender, req.Message)))

	msg := &pb.ChatMessage{Sender: req.Sender, Message: req.Message, Timestamp: uint64(time.Now().Unix())}
	err := publishMessage(s, msg)

	if err != nil {
		s.Logger.Error("failed to publish message", err)
		return &pb.SendChatMessageResponse{Sender: msg.Sender, Timestamp: msg.Timestamp, Success: false}, nil
	} else {
		s.Logger.Info(fmt.Sprintf("[#%d] (jetstream) Published message to [chat.global]: %s", s.sentMessages+1, req.Message))
	}

	s.sentMessages++
	return &pb.SendChatMessageResponse{Sender: msg.Sender, Timestamp: msg.Timestamp, Success: true}, nil
}

func InitNats() (*nats.Conn, error) {
	// nats
	logger := logging.Get("chatserver", false)
	cfg, err := config.LoadConfig()
	logger.Info("loaded config", "nats_url", cfg.NatsURL, "max_async_publish", cfg.MaxAsyncPublish)

	if err != nil {
		logger.Error("failed to load config", err)
		return nil, err
	}

	nc, err := nats.Connect(cfg.NatsURL)

	if err != nil {
		logger.Error("failed to connect to nats", err)
		return nil, err
	}

	// init nats-js stream
	js, err := nc.JetStream(nats.PublishAsyncMaxPending(cfg.MaxAsyncPublish))
	if err != nil {
		logger.Error("connection to jetstream failed", err)
		return nil, err
	}

	_, err = js.StreamInfo("chat")
	if err != nil {
		logger.Warn("failed to get stream info", err)

		if err.Error() == nats.ErrStreamNotFound.Error() {
			js.AddStream(&nats.StreamConfig{Name: "chat", Subjects: []string{"chat.global"}, Storage: nats.FileStorage})
			logger.Info("created stream", "name", "chat")
			return nc, nil
		}

		return nil, err
	} else {
		logger.Info("stream already exists, skipping create", "name", "chat")
	}

	return nc, nil
}

package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/log"
	pb "github.com/komadiina/spelltext/proto/chat"
	"github.com/komadiina/spelltext/utils/singleton/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr     = flag.String("addr", "localhost:50051", "server address")
	username = flag.String("username", "John_Doe", "name to identify sent messages")
)

func SendMessage(content string, client pb.ChatServiceClient, ctx context.Context) {
	logger := logging.Get()

	resp, err := client.SendChatMessage(ctx, &pb.SendChatMessageRequest{Sender: *username, Message: content})
	if err != nil {
		logger.Error("failed to send message", err)
		os.Exit(1)
	}

	logger.Info(fmt.Sprintf("Received response from %s.", resp.GetSender()), "success", resp.GetSuccess())
}

func main() {
	logging.Init(log.InfoLevel)
	logger := logging.Get()

	flag.Parse()

	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error("failed to dial server", err)
		os.Exit(1)
	}
	defer conn.Close()
	client := pb.NewChatServiceClient(conn)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		msg, err := reader.ReadString('\n')
		if err != nil {
			logger.Error("failed to read message", err)
		}

		msg = strings.Trim(msg, "\r\n")

		if msg == "/exit" {
			break
		}

		SendMessage(msg, client, ctx)
	}
}

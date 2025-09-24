package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	pb "github.com/komadiina/spelltext/proto/chat"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var logger log.Logger

var (
	addr = flag.String("addr", "127.0.0.1:50051", "server address")
)

func SendMessage(content string, client pb.ChatServiceClient, ctx context.Context) {
	resp, err := client.SendChatMessage(ctx, &pb.SendChatMessageRequest{Sender: "Bob", Message: content})
	if err != nil {
		logger.Error("failed to send message", err)
		os.Exit(1)
	}

	logger.Info(fmt.Sprintf("Received response from %s.", resp.GetSender()), "success", resp.GetSuccess())
}

func main() {
	logger = *log.NewWithOptions(os.Stdout, log.Options{
		ReportCaller:    true,
		ReportTimestamp: true,
		TimeFormat:      time.DateTime,
	})
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

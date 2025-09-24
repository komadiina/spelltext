package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/gdamore/tcell/v2"
	pb "github.com/komadiina/spelltext/proto/chat"
	"github.com/komadiina/spelltext/utils/singleton/logging"
	"github.com/rivo/tview"
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

func JoinChatroom(ctx context.Context, client pb.ChatServiceClient) {
	client.JoinChatroom(ctx, &pb.JoinChatroomMessageRequest{Username: *username})
}

func LeaveChatroom(ctx context.Context, client pb.ChatServiceClient) {
	client.LeaveChatroom(ctx, &pb.LeaveChatroomMessageRequest{Username: *username})
}

func main() {
	logging.Init(log.ErrorLevel)
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
	JoinChatroom(ctx, client)
	defer LeaveChatroom(ctx, client)
	defer cancel()

	app := tview.NewApplication()
	chat := tview.NewTextView().
		SetTextAlign(tview.AlignLeft).
		SetRegions(false).
		ScrollToEnd()

	go func() {
		stream, err := client.Subscribe(ctx, &pb.SubscribeRequest{Username: *username})
		if err != nil {
			logger.Error("failed to subscribe", err)
			os.Exit(1)
		}

		for {
			msg, err := stream.Recv()
			if err != nil || err == io.EOF {
				logger.Error("failed to receive message", err)
				os.Exit(1)
			}

			app.QueueUpdateDraw(func() {
				fmt.Fprintf(chat, "[%s]: %s\n", msg.Sender, msg.Message)
				chat.ScrollToEnd()
			})
		}
	}()

	input := tview.NewInputField().
		SetLabel("> ")

	input.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter && strings.Trim(input.GetText(), "\r\n") != "" {
			SendMessage(input.GetText(), client, ctx)
			input.SetText("")
		} else {
			input.SetText("")
		}
	})

	chat.SetBorder(true).SetTitle(" chat ")
	input.SetBorder(true).SetTitle(" input ")

	grid := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(chat, 0, 4, false).
		AddItem(input, 3, 1, true)

	grid.SetBorder(true).SetBorderPadding(1, 0, 1, 0).SetTitle(" spelltext - chat ")

	if err := app.SetRoot(grid, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

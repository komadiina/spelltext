package views

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/charmbracelet/log"
	"github.com/gdamore/tcell/v2"
	"github.com/komadiina/spelltext/client/config"
	types "github.com/komadiina/spelltext/client/types"
	pb "github.com/komadiina/spelltext/proto/chat"
	"github.com/komadiina/spelltext/utils/singleton/logging"
	"github.com/nats-io/nats.go"
	"github.com/rivo/tview"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
)

var (
	addr     = flag.String("addr", "localhost:50051", "server address")
	username = flag.String("username", "John_Doe", "name to identify sent messages")
)

func SendMessage(content string, client pb.ChatServiceClient, ctx context.Context) {
	logger := logging.Get("chatserver")

	resp, err := client.SendChatMessage(ctx, &pb.SendChatMessageRequest{Sender: *username, Message: content})
	if err != nil {
		logger.Error("failed to send message", "reason", err)
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

func InitJetStream(cfg *config.Config) (*nats.Conn, nats.JetStream, error) {
	conn, err := nats.Connect(cfg.NatsURL)
	if err != nil {
		return nil, nil, err
	}

	js, err := conn.JetStream()
	if err != nil {
		return conn, nil, err
	}

	return conn, js, nil
}

func AddChatPage(c *types.SpelltextClient) *tview.Pages {
	return nil
}

func main() {
	logging.Init(log.ErrorLevel, "client")
	logger := logging.Get("client")

	flag.Parse()

	cfg, err := config.LoadConfig()

	app := tview.NewApplication()

	nc, js, err := InitJetStream(cfg)
	if err != nil {
		logger.Error("failed to init jetstream", "reason", err)
		os.Exit(1)
	}
	defer nc.Drain()

	chat := tview.NewTextView().
		SetTextAlign(tview.AlignLeft).
		SetRegions(false).
		ScrollToEnd()
	sub, err := js.Subscribe("chat.global", func(msg *nats.Msg) {
		var chatMsg pb.ChatMessage
		err := proto.Unmarshal(msg.Data, &chatMsg)
		if err != nil {
			logger.Error("failed to unmarshal message", "reason", err)
			return
		}
		logger.Info(fmt.Sprintf("Received message from %s: %s", chatMsg.Sender, chatMsg.Message))

		app.QueueUpdateDraw(func() {
			fmt.Fprintf(chat, fmt.Sprint(chatMsg.Message))
			chat.ScrollToEnd()
		})
	})
	if err != nil {
		logger.Error("failed to subscribe to 'chat.global'", "reason", err)
		os.Exit(1)
	}
	defer sub.Drain()

	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error("failed to dial server", "reason", err)
		os.Exit(1)
	}
	defer conn.Close()

	logger.Info("connected to server!", "addr", *addr)
	logger.Info("initializing nats jetstream subscription...")

	client := pb.NewChatServiceClient(conn)
	ctx, ctxCancel := context.WithCancel(context.Background())
	JoinChatroom(ctx, client)
	defer LeaveChatroom(ctx, client)
	defer ctxCancel()

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

	// graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	wg := sync.WaitGroup{}
	wg.Go(func() {
		_ = <-sigCh
		LeaveChatroom(ctx, client)
		ctxCancel()
	})

	if err := app.SetRoot(grid, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

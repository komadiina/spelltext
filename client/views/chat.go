package views

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/komadiina/spelltext/client/constants"
	types "github.com/komadiina/spelltext/client/types"
	pb "github.com/komadiina/spelltext/proto/chat"
	"github.com/nats-io/nats.go"
	"github.com/rivo/tview"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
)

func SendMessage(content string, client pb.ChatClient, ctx context.Context, c *types.SpelltextClient) {
	resp, err := client.SendChatMessage(
		ctx, &pb.SendChatMessageRequest{
			Sender:  c.User.Username,
			Message: content,
		})

	if err != nil {
		c.Logger.Error("failed to send message", "reason", err)
		os.Exit(1)
	}

	c.Logger.Info(fmt.Sprintf("Received response from %s.", resp.GetSender()), "success", resp.GetSuccess())
}

func JoinChatroom(ctx context.Context, client pb.ChatClient, c *types.SpelltextClient) {
	client.JoinChatroom(ctx, &pb.JoinChatroomMessageRequest{Username: c.User.Username})
}

func LeaveChatroom(ctx context.Context, client pb.ChatClient, c *types.SpelltextClient) {
	client.LeaveChatroom(ctx, &pb.LeaveChatroomMessageRequest{Username: c.User.Username})
}

func AddChatPage(c *types.SpelltextClient) {
	onClose := func() {}

	c.PageManager.RegisterFactory(constants.PAGE_CHAT, func() tview.Primitive {
		chat := tview.NewTextView().
			SetTextAlign(tview.AlignLeft).
			SetRegions(false).
			ScrollToEnd()

		js, err := c.Nats.JetStream()
		if err != nil {
			c.Logger.Error("failed to get jetstream", "reason", err)
			c.PageManager.Pop()
		}

		_, err = js.Subscribe("chat.global", func(msg *nats.Msg) {
			var chatMsg pb.ChatMessage
			err := proto.Unmarshal(msg.Data, &chatMsg)
			if err != nil {
				c.Logger.Error("failed to unmarshal message", "reason", err)
				return
			}
			c.Logger.Info(fmt.Sprintf("Received message from %s: %s", chatMsg.Sender, chatMsg.Message))

			c.App.QueueUpdateDraw(func() {
				fmt.Fprintf(chat, fmt.Sprint(chatMsg.Message))
				chat.ScrollToEnd()
			})
		})
		if err != nil {
			c.Logger.Error("failed to subscribe to 'chat.global'", "reason", err)
			os.Exit(1)
		}

		conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			c.Logger.Error("failed to dial server", "reason", err)
			os.Exit(1)
		}

		c.Logger.Info("connected to server!", "addr", conn.Target())
		c.Logger.Info("initializing nats jetstream subscription...")

		client := pb.NewChatClient(conn)
		ctx, cancelFunc := context.WithCancel(context.Background())
		JoinChatroom(ctx, client, c)

		input := tview.NewInputField().
			SetLabel("> ")

		input.SetDoneFunc(func(key tcell.Key) {
			sanitized := strings.Trim(input.GetText(), "\r\n")

			if key == tcell.KeyEnter && sanitized != "" {
				if sanitized == "/exit" {
					c.PageManager.Pop()
					return
				}

				SendMessage(input.GetText(), client, ctx, c)
				input.SetText("")
			} else {
				input.SetText("")
			}
		})

		chat.SetBorder(true).SetTitle(" chat ")
		input.SetBorder(true).SetTitle(" input ")

		flex := tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(chat, 0, 4, false).
			AddItem(input, 3, 1, true).
			SetFullScreen(true)

		// register cleanup method
		onClose = func() {
			fmt.Println("onClose")
			LeaveChatroom(ctx, client, c)
			cancelFunc()
			conn.Close()
		}

		return flex
	}, nil, onClose)
}

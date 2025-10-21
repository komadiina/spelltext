package views

import (
	"fmt"
	"os"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/komadiina/spelltext/client/constants"
	types "github.com/komadiina/spelltext/client/types"
	pb "github.com/komadiina/spelltext/proto/chat"
	"github.com/nats-io/nats.go"
	"github.com/rivo/tview"
	"google.golang.org/protobuf/proto"
)

func SendMessage(content string, c *types.SpelltextClient) {
	resp, err := c.Clients.ChatClient.SendChatMessage(
		*c.Context, &pb.SendChatMessageRequest{
			Sender:  c.User.Username,
			Message: content,
		})

	if err != nil {
		c.Logger.Error("failed to send message", "reason", err)
		os.Exit(1)
	}

	c.Logger.Info(fmt.Sprintf("Received response from %s.", resp.GetSender()), "success", resp.GetSuccess())
}

func JoinChatroom(c *types.SpelltextClient) {
	c.Clients.ChatClient.JoinChatroom(*c.Context, &pb.JoinChatroomMessageRequest{Username: c.User.Username})
}

func LeaveChatroom(c *types.SpelltextClient) {
	c.Clients.ChatClient.LeaveChatroom(*c.Context, &pb.LeaveChatroomMessageRequest{Username: c.User.Username})
}

func AddChatPage(c *types.SpelltextClient) {
	c.PageManager.RegisterFactory(constants.PAGE_CHAT, func() tview.Primitive {
		chat := tview.NewTextView().
			SetTextAlign(tview.AlignLeft).
			ScrollToEnd()
		chat.SetDynamicColors(true)

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
				fmt.Fprintf(chat, `[%s]%s[""]: %s%s`, "#EBDBB2", tview.Escape("["+chatMsg.Sender+"]"), chatMsg.Message, "\n")
				chat.ScrollToEnd()
			})
		})
		if err != nil {
			c.Logger.Error("failed to subscribe to 'chat.global'", "reason", err)
			os.Exit(1)
		}

		JoinChatroom(c)

		input := tview.NewInputField().
			SetLabel("> ")

		input.SetDoneFunc(func(key tcell.Key) {
			sanitized := strings.Trim(input.GetText(), "\r\n")

			if key == tcell.KeyEnter && sanitized != "" {
				switch sanitized {
				case "/exit":
					c.PageManager.Pop()
					return
				case "/clear":
					chat.Clear()
					input.SetText("")
					return
				}

				SendMessage(input.GetText(), c)
				input.SetText("")
			} else {
				input.SetText("")
			}
		})

		chat.SetBorder(true).SetTitle(" [::b]chat[::-] ")
		input.SetBorder(true).SetTitle(" [::b]input[::-] ")

		flex := tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(chat, 0, 4, false).
			AddItem(input, 3, 1, true).
			SetFullScreen(true)

		return flex
	}, nil, func() { LeaveChatroom(c) })
}

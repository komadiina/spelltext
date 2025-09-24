package pages

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/rivo/tview"
)

var users = []string{"Bob", "Alice"}
var messages = []string{"Hello", "World", "Foo", "Bar"}

type ChatMessage struct {
	User    string
	Message string
}

func GenerateChat(app *tview.Application, pages *tview.Pages) *tview.Pages {
	chat := tview.NewTextView().
		SetTextAlign(tview.AlignLeft).
		SetRegions(false).
		ScrollToEnd()

	const numMessages = 50
	go func() {
		for i := 0; i < numMessages; i++ {
			time.Sleep(time.Second)
			var message = generateMessage()
			escaped := fmt.Sprintf("[%s]: %s\n", message.User, message.Message)

			app.QueueUpdateDraw(func() {
				fmt.Fprint(chat, escaped)
			})
		}
	}()

	chat.SetBorder(true).SetTitle(" chat ")

	pages.AddPage("chat", chat, true, true)
	return pages
}

func generateMessage() ChatMessage {
	user := users[rand.Intn(len(users))]
	message := messages[rand.Intn(len(messages))]
	return ChatMessage{User: user, Message: message}
}

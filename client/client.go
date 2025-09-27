package main

import (
	"flag"

	"github.com/charmbracelet/log"
	"github.com/gdamore/tcell/v2"
	"github.com/komadiina/spelltext/client/config"
	"github.com/komadiina/spelltext/client/factory"
	"github.com/komadiina/spelltext/client/registry"
	"github.com/komadiina/spelltext/client/types"
	"github.com/komadiina/spelltext/client/views"
	"github.com/komadiina/spelltext/utils/singleton/logging"
	"github.com/nats-io/nats.go"
	"github.com/rivo/tview"
)

func InitNats(cfg *config.Config) (*nats.Conn, nats.JetStream, error) {
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

func InitRegistry(cfg *config.Config) *registry.Registry {
	return registry.NewRegistry()
}

func InitializePages(client *types.SpelltextClient) {
	client.PageManager.Pages.
		SetBorder(true).
		SetTitle(`[blueviolet]╝[""] [white]spelltext[""] [blueviolet]╚[""]`).
		SetBorderPadding(2, 2, 10, 10).
		SetBorderStyle(tcell.StyleDefault.Foreground(tcell.ColorBlueViolet))

	views.AddLoginPage(client)
	views.AddMainmenuPage(client)
	views.AddChatPage(client)
	views.AddCharactersPage(client)
	views.AddStorePage(client)
	views.AddProgressPage(client)
	views.AddGambaPage(client)
	views.AddCombatPage(client)
	views.AddInventoryPage(client)
}

func main() {
	flag.Parse()
	logger := logging.Get("client")
	logger.SetLevel(log.ErrorLevel)

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("failed to load config", "reason", err)
	}

	client := types.SpelltextClient{
		Config: cfg,
		Logger: logger,
		App:    tview.NewApplication(),
		User:   &types.SpelltextUser{},
	}

	nc, _, err := InitNats(cfg)
	if err != nil {
		logger.Fatal("failed to init nats/jetstream", "reason", err)
	}

	client.Nats = nc
	client.PageManager = factory.NewPageManager(client.App)
	InitializePages(&client)

	client.NavigateTo = func(pageKey string) {
		if client.PageManager.HasPage(pageKey) == false {
			logger.Fatal("page not found", "page", pageKey)
			return
		}

		client.PageManager.Push(pageKey, false)
	}

	client.NavigateTo(views.LOGIN_PAGE)

	client.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			client.PageManager.Pop()
		}

		return event
	})

	if err := client.App.SetRoot(client.PageManager.Pages, true).EnableMouse(true).Run(); err != nil {
		logger.Fatal(err)
	}

	// cleanup
	client.Nats.Drain()
}

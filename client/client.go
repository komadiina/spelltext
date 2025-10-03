package main

import (
	"context"
	"flag"
	"fmt"
	"os"

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
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pbChat "github.com/komadiina/spelltext/proto/chat"
	pbStore "github.com/komadiina/spelltext/proto/store"
)

func InitializeNats(cfg *config.Config) (*nats.Conn, nats.JetStream, error) {
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
	views.AddLoginPage(client)
	views.AddMainmenuPage(client)
	views.AddChatPage(client)
	views.AddCharactersPage(client)
	views.AddStorePage(client)
	views.AddProgressPage(client)
	views.AddGambaPage(client)
	views.AddCombatPage(client)
	views.AddInventoryPage(client)
	views.AddVendorPage(client)
}

func InitializeClients(c *types.SpelltextClient) {
	c.Clients = &types.Clients{}

	// init chat client
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.Logger.Error("failed to init chat client", "reason", err)
	} else {
		c.Clients.ChatClient = pbChat.NewChatClient(conn)
	}

	conn, err = grpc.NewClient("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.Logger.Error("failed to init store client", "reason", err)
	} else {
		c.Clients.StoreClient = pbStore.NewStoreClient(conn)
	}
}

func main() {
	flag.Parse()
	logging.Init(log.InfoLevel, "client", true)
	logger := logging.Get("client", true)

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("failed to load config", "reason", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	client := types.SpelltextClient{
		Config:     cfg,
		Logger:     logger,
		App:        tview.NewApplication(),
		User:       &types.SpelltextUser{},
		Context:    &ctx,
		AppStorage: make(map[string]any),
	}
	InitializeClients(&client)

	nc, _, err := InitializeNats(cfg)
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

	client.NavigateTo(views.PAGE_LOGIN)

	client.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			client.PageManager.Pop()
		}

		return event
	})

	if err := client.App.SetRoot(client.PageManager.Pages, true).EnableMouse(true).Run(); err != nil {
		client.App.Stop()
		client.Logger.Error(err)
		panic(err)
	}

	fmt.Fprint(os.Stderr, "\x1b[?1049l") // switch back to main screen
	// then print error or call external `reset` if needed
	fmt.Fprintln(os.Stderr, "app.Run error:", err)

	// cleanup
	client.Nats.Drain()
	defer cancel()
}

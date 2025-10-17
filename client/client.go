package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/gdamore/tcell/v2"
	"github.com/komadiina/spelltext/client/audio"
	"github.com/komadiina/spelltext/client/config"
	"github.com/komadiina/spelltext/client/constants"
	"github.com/komadiina/spelltext/client/factory"
	"github.com/komadiina/spelltext/client/types"
	"github.com/komadiina/spelltext/client/views"
	"github.com/komadiina/spelltext/utils/singleton/logging"
	"github.com/nats-io/nats.go"
	"github.com/rivo/tview"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pbArmory "github.com/komadiina/spelltext/proto/armory"
	pbChat "github.com/komadiina/spelltext/proto/chat"
	pbGamba "github.com/komadiina/spelltext/proto/gamba"
	pbInventory "github.com/komadiina/spelltext/proto/inventory"
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


func InitializePages(client *types.SpelltextClient) {
	views.AddLoginPage(client)
	views.AddMainmenuPage(client)
	views.AddChatPage(client)
	views.AddArmoryPage(client)
	views.AddStorePage(client)
	views.AddProgressPage(client)
	views.AddGambaPage(client)
	views.AddCombatPage(client)
	views.AddInventoryPage(client)
	views.AddVendorPage(client)
}

func InitializeClients(c *types.SpelltextClient) {
	c.Clients = &types.Clients{}

	// chat
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.Logger.Error("failed to init chat client", "reason", err)
	} else {
		c.Clients.ChatClient = pbChat.NewChatClient(conn)
	}

	// store
	conn, err = grpc.NewClient("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.Logger.Error("failed to init store client", "reason", err)
	} else {
		c.Clients.StoreClient = pbStore.NewStoreClient(conn)
	}

	// inventory
	conn, err = grpc.NewClient("localhost:50053", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.Logger.Error("failed to init inventory client", "reason", err)
	} else {
		c.Clients.InventoryClient = pbInventory.NewInventoryClient(conn)
	}

	// armory
	conn, err = grpc.NewClient("localhost:50054", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.Logger.Error("failed to init armory client", "reason", err)
	} else {
		c.Clients.CharacterClient = pbArmory.NewCharacterClient(conn)
	}

	// gamba
	conn, err = grpc.NewClient("localhost:50055", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.Logger.Error("failed to init gamba client", "reason", err)
	} else {
		c.Clients.GambaClient = pbGamba.NewGambaClient(conn)
	}
}

func main() {
	flag.Parse()
	logging.Init(log.InfoLevel, "client", true)
	logger := logging.Get("client", true)

	logger.Info("loading client config...", "CONFIG_FILE", os.Getenv("CONFIG_FILE"))
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error("failed to load config, using default values.", "reason", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	// init audio manager
	am := audio.NewManager(44000)
	err = am.Preload(constants.PRELOAD)

	if err != nil {
		logger.Fatal("failed to preload audio files", "reason", err)
	}

	client := types.SpelltextClient{
		Config:       cfg,
		Logger:       logger,
		App:          tview.NewApplication(),
		User:         &types.SpelltextUser{},
		Context:      &ctx,
		AppStorage:   make(map[string]any),
		AudioManager: am,
	}

	logger.Info("initializing clients...")
	InitializeClients(&client)
	logger.Info("clients initialized.")

	logger.Info("initializng nats...")
	nc, _, err := InitializeNats(cfg)
	if err != nil {
		logger.Fatal("failed to init nats/jetstream", "reason", err)
	}

	client.Nats = nc
	logger.Info("nats/js initialized.")

	logger.Info("initializing PageManager factory...")
	client.PageManager = factory.NewPageManager(client.Logger, client.App)
	InitializePages(&client)
	logger.Info("PageManager factory initialized.")

	client.NavigateTo = func(pageKey string) {
		if client.PageManager.HasPage(pageKey) == false {
			logger.Fatal("page not found", "page", pageKey)
			return
		}

		client.PageManager.Push(pageKey, false)
	}

	client.NavigateTo(constants.PAGE_LOGIN)

	client.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			client.AudioManager.Play(constants.BLIP_TINY, client.Logger)

			if client.PageManager.Pop() == -1 {
				client.App.Stop()
				return nil
			}
		} else if event.Key() == tcell.KeyEnter {
			client.AudioManager.Play(constants.BLIP_NOTIFICATION, client.Logger)
		} else {
			client.AudioManager.Play(constants.BLIP_NAVIGATE, client.Logger)
		}

		return event
	})

	if err := client.App.SetRoot(client.PageManager.Pages, true).EnableMouse(true).Run(); err != nil {
		client.Logger.Error(err)
	}

	// cleanup
	client.Nats.Drain()
	defer cancel()

	logger.Info("client shutdown.")
	fmt.Println("bye")
}

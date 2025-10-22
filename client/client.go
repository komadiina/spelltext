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

	pbAuth "github.com/komadiina/spelltext/proto/auth"
	pbChar "github.com/komadiina/spelltext/proto/char"
	pbChat "github.com/komadiina/spelltext/proto/chat"
	pbCombat "github.com/komadiina/spelltext/proto/combat"
	pbGamba "github.com/komadiina/spelltext/proto/gamba"
	pbInventory "github.com/komadiina/spelltext/proto/inventory"
	pbRepo "github.com/komadiina/spelltext/proto/repo"
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
	views.AddCharacterPage(client)
	views.AddStorePage(client)
	views.AddProgressPage(client)
	views.AddGambaPage(client)
	views.AddCombatPage(client)
	views.AddInventoryPage(client)
	views.AddVendorPage(client)
	views.AddCreateCharacterPage(client)
	views.AddFightPage(client)
}

func InitializeClients(c *types.SpelltextClient) {
	c.Clients = &types.Clients{}
	c.Connections = &types.Connections{}

	host := func(port int) string { return fmt.Sprintf("localhost:%d", port) }

	// chat
	conn, err := grpc.NewClient(host(c.Config.ChatPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.Logger.Error("failed to init chat client", "reason", err)
	} else {
		c.Clients.ChatClient = pbChat.NewChatClient(conn)
		c.Connections.Chat = conn
	}

	// store
	conn, err = grpc.NewClient(host(c.Config.StorePort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.Logger.Error("failed to init store client", "reason", err)
	} else {
		c.Clients.StoreClient = pbStore.NewStoreClient(conn)
		c.Connections.Store = conn
	}

	// inventory
	conn, err = grpc.NewClient(host(c.Config.InventoryPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.Logger.Error("failed to init inventory client", "reason", err)
	} else {
		c.Clients.InventoryClient = pbInventory.NewInventoryClient(conn)
		c.Connections.Inventory = conn
	}

	// character
	conn, err = grpc.NewClient(host(c.Config.CharacterPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.Logger.Error("failed to init character client", "reason", err)
	} else {
		c.Clients.CharacterClient = pbChar.NewCharacterClient(conn)
		c.Connections.Character = conn
	}

	// gamba
	conn, err = grpc.NewClient(host(c.Config.GambaPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.Logger.Error("failed to init gamba client", "reason", err)
	} else {
		c.Clients.GambaClient = pbGamba.NewGambaClient(conn)
		c.Connections.Gamba = conn
	}

	// auth
	conn, err = grpc.NewClient(host(c.Config.AuthPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.Logger.Error("failed to init auth client", "reason", err)
	} else {
		c.Clients.AuthClient = pbAuth.NewAuthClient(conn)
		c.Connections.Auth = conn
	}

	// combat
	conn, err = grpc.NewClient(host(c.Config.CombatPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.Logger.Error("failed to init auth client", "reason", err)
	} else {
		c.Clients.CombatClient = pbCombat.NewCombatClient(conn)
		c.Connections.Combat = conn
	}
}

func closeClients(c *types.SpelltextClient) {
	if c.Clients != nil {
		c.Connections.Chat.Close()
		c.Connections.Store.Close()
		c.Connections.Inventory.Close()
		c.Connections.Character.Close()
		c.Connections.Gamba.Close()
		c.Connections.Auth.Close()
		c.Connections.Combat.Close()
	}
}

var (
	fDebugLevel = flag.String("debug", "info", "debug level")
)

const banner = `

                _ _ _            _   
               | | | |          | |  
 ___ _ __   ___| | | |_ _____  _| |_ 
/ __| '_ \ / _ \ | | __/ _ \ \/ / __|
\__ \ |_) |  __/ | | ||  __/>  <| |_ 
|___/ .__/ \___|_|_|\__\___/_/\_\\__|
    | |                              
    |_|                              

`

func main() {
	flag.Parse()
	var level log.Level = log.InfoLevel
	if fDebugLevel == nil {
		level = log.InfoLevel
	} else {
		level, _ = log.ParseLevel(*fDebugLevel)
	}

	logging.Init(level, "client", true)
	logger := logging.Get("client", true)

	err := os.Setenv("CONFIG_FILE", "config.yml")
	logger.Debug("loading client config...", "CONFIG_FILE", os.Getenv("CONFIG_FILE"))
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

	am.AudioEnabled = cfg.AudioEnabled

	client := types.SpelltextClient{
		Config:       cfg,
		Logger:       logger,
		App:          tview.NewApplication(),
		User:         &types.SpelltextUser{},
		Context:      &ctx,
		AppStorage:   make(map[string]any),
		AudioManager: am,
		Storage: &types.AppStorage{
			Ministate:         &types.Ministate{},
			CurrentUser:       &pbRepo.User{},
			SelectedCharacter: &pbRepo.Character{},
			SelectedVendor:    &pbRepo.Vendor{},
			EquipSlots:        nil,
		},
	}

	client.AudioManager.PlayBackground(logger)

	logger.Debug("initializing clients...")
	InitializeClients(&client)
	defer closeClients(&client)
	logger.Debug("clients initialized.")

	logger.Debug("initializng nats...")
	nc, _, err := InitializeNats(cfg)
	if err != nil {
		logger.Fatal("failed to init nats/jetstream", "reason", err)
	}

	client.Nats = nc
	logger.Debug("nats/js initialized.")

	logger.Debug("initializing PageManager factory...")
	client.PageManager = factory.NewPageManager(client.Logger, client.App)
	InitializePages(&client)
	logger.Debug("PageManager factory initialized.")

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
			client.AudioManager.Play(constants.BLIP_BACKWARD, client.Logger)

			if client.PageManager.Pop() == -1 {
				client.App.Stop()
				return nil
			}
		} else if event.Key() == tcell.KeyEnter {
			client.AudioManager.Play(constants.BLIP_FORWARD, client.Logger)
		} else {
			client.AudioManager.Play(constants.BLIP_INPUT, client.Logger)
		}

		return event
	})

	if err := client.App.SetRoot(client.PageManager.Pages, true).EnableMouse(true).Run(); err != nil {
		client.Logger.Error(err)
	}

	// cleanup
	client.Nats.Drain()
	defer cancel()
	logger.Debug("client shutdown.")

	fmt.Print(banner)
	goodbye := `
> thanks for playing this torturefest
> a game by ogg/komadiina (https://github.com/komadiina)
> follow the development at https://github.com/komadiina/spelltext
~ kthxb
`
	fmt.Print(goodbye)
}

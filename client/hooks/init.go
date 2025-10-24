package hooks

import (
	"fmt"

	"github.com/komadiina/spelltext/client/config"
	"github.com/komadiina/spelltext/client/types"
	"github.com/komadiina/spelltext/client/views"
	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pbAuth "github.com/komadiina/spelltext/proto/auth"
	pbBuild "github.com/komadiina/spelltext/proto/build"
	pbChar "github.com/komadiina/spelltext/proto/char"
	pbChat "github.com/komadiina/spelltext/proto/chat"
	pbCombat "github.com/komadiina/spelltext/proto/combat"
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
	views.AddCharacterPage(client)
	views.AddStorePage(client)
	views.AddProgressPage(client)
	views.AddGambaPage(client)
	views.AddCombatPage(client)
	views.AddInventoryPage(client)
	views.AddVendorPage(client)
	views.AddCreateCharacterPage(client)
	views.AddFightPage(client)
	views.AddAbilityPage(client)
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

	// build
	conn, err = grpc.NewClient(host(c.Config.BuildPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.Logger.Error("failed to init build client", "reason", err)
	} else {
		c.Clients.BuildClient = pbBuild.NewBuildClient(conn)
		c.Connections.Build = conn
	}
}

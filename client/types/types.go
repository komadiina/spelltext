package types

import (
	"context"

	"github.com/komadiina/spelltext/client/audio"
	"github.com/komadiina/spelltext/client/config"
	"github.com/komadiina/spelltext/client/factory"
	pbAuth "github.com/komadiina/spelltext/proto/auth"
	pbBuild "github.com/komadiina/spelltext/proto/build"
	pbChar "github.com/komadiina/spelltext/proto/char"
	pbChat "github.com/komadiina/spelltext/proto/chat"
	pbCombat "github.com/komadiina/spelltext/proto/combat"
	pbGamba "github.com/komadiina/spelltext/proto/gamba"
	pbInventory "github.com/komadiina/spelltext/proto/inventory"
	pbRepo "github.com/komadiina/spelltext/proto/repo"
	pbStore "github.com/komadiina/spelltext/proto/store"
	"github.com/komadiina/spelltext/utils/singleton/logging"
	"github.com/nats-io/nats.go"
	"github.com/rivo/tview"
	"google.golang.org/grpc"
)

type SpelltextClient struct {
	Config       *config.Config
	Logger       *logging.Logger
	Nats         *nats.Conn
	App          *tview.Application
	AppStorage   map[string]any
	Storage      *AppStorage
	PageManager  *factory.PageManager
	Clients      *Clients
	Connections  *Connections
	Context      *context.Context
	AudioManager *audio.Manager

	NavigateTo func(pageKey string)
}

type Ministate struct {
	Username   string
	CurrentNpc *pbRepo.Npc
	FightState *NpcFightState
}

type AppStorage struct {
	Ministate         *Ministate
	CurrentUser       *pbRepo.User
	SelectedCharacter *pbRepo.Character
	SelectedVendor    *pbRepo.Vendor
	EquipSlots        []*pbRepo.EquipSlot
	CharacterStats    *CharacterStats
}

type Clients struct {
	ChatClient      pbChat.ChatClient
	StoreClient     pbStore.StoreClient
	CharacterClient pbChar.CharacterClient
	InventoryClient pbInventory.InventoryClient
	GambaClient     pbGamba.GambaClient
	AuthClient      pbAuth.AuthClient
	CombatClient    pbCombat.CombatClient
	BuildClient     pbBuild.BuildClient
}

type Connections struct {
	Inventory *grpc.ClientConn
	Character *grpc.ClientConn
	Chat      *grpc.ClientConn
	Store     *grpc.ClientConn
	Gamba     *grpc.ClientConn
	Auth      *grpc.ClientConn
	Combat    *grpc.ClientConn
	Build     *grpc.ClientConn
}

type UnusableHotkey struct {
	Key  string
	Desc string
}

type CharacterStats struct {
	HealthPoints int64
	PowerPoints  int64
	Strength     int64
	Spellpower   int64
	Armor        int64
	Damage       int64
}

type EntityStatusFrame struct {
	Health    uint64
	Power     uint64
	BarHealth *tview.TextView
	BarPower  *tview.TextView
	FlHealth  *tview.Flex
	FlPower   *tview.Flex
	FlTextual *tview.Flex
	Refresh   func(newHp int, newPwr int)
}

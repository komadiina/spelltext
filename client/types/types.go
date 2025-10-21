package types

import (
	"context"

	"github.com/komadiina/spelltext/client/audio"
	"github.com/komadiina/spelltext/client/config"
	"github.com/komadiina/spelltext/client/factory"
	pbAuth "github.com/komadiina/spelltext/proto/auth"
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
	User         *SpelltextUser
	Clients      *Clients
	Connections  *Connections
	Context      *context.Context
	AudioManager *audio.Manager

	NavigateTo func(pageKey string)
}

type Ministate struct {
	Username string
}

type AppStorage struct {
	Ministate         *Ministate
	CurrentUser       *pbRepo.User
	SelectedCharacter *pbRepo.Character
	SelectedVendor    *pbRepo.Vendor
	EquipSlots        []*pbRepo.EquipSlot
}

type SpelltextUser struct {
	Username string
}

type ContextDef struct {
	Context context.Context
	Cancel  context.CancelFunc
}

type Clients struct {
	ChatClient      pbChat.ChatClient
	StoreClient     pbStore.StoreClient
	CharacterClient pbChar.CharacterClient
	InventoryClient pbInventory.InventoryClient
	GambaClient     pbGamba.GambaClient
	AuthClient      pbAuth.AuthClient
	CombatClient    pbCombat.CombatClient
}

type Connections struct {
	Inventory *grpc.ClientConn
	Character *grpc.ClientConn
	Chat      *grpc.ClientConn
	Store     *grpc.ClientConn
	Gamba     *grpc.ClientConn
	Auth      *grpc.ClientConn
	Combat    *grpc.ClientConn
}

type NavigableForm struct {
	tview.Form
}

type NavigableFormButton struct {
	tview.Button

	LeftNbr   *tview.Button
	TopNbr    *tview.Button
	RightNbr  *tview.Button
	BottomNbr *tview.Button
}

type CharacterStats struct {
	HealthPoints int64
	PowerPoints  int64
	Strength     int64
	Spellpower   int64
	Armor        int64
	Damage       int64
}

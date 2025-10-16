package types

import (
	"context"

	"github.com/komadiina/spelltext/client/audio"
	"github.com/komadiina/spelltext/client/config"
	"github.com/komadiina/spelltext/client/factory"
	pbArmory "github.com/komadiina/spelltext/proto/armory"
	pbChat "github.com/komadiina/spelltext/proto/chat"
	pbInventory "github.com/komadiina/spelltext/proto/inventory"
	pbStore "github.com/komadiina/spelltext/proto/store"
	"github.com/komadiina/spelltext/utils/singleton/logging"
	"github.com/nats-io/nats.go"
	"github.com/rivo/tview"
)

type SpelltextClient struct {
	Config       *config.Config
	Logger       *logging.Logger
	Nats         *nats.Conn
	App          *tview.Application
	AppStorage   map[string]any
	PageManager  *factory.PageManager
	User         *SpelltextUser
	Clients      *Clients
	Context      *context.Context
	AudioManager *audio.Manager

	NavigateTo func(pageKey string)
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
	CharacterClient pbArmory.CharacterClient
	InventoryClient pbInventory.InventoryClient
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

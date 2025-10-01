package types

import (
	"context"

	"github.com/komadiina/spelltext/client/config"
	"github.com/komadiina/spelltext/client/factory"
	pbChat "github.com/komadiina/spelltext/proto/chat"
	pbStore "github.com/komadiina/spelltext/proto/store"
	"github.com/komadiina/spelltext/utils/singleton/logging"
	"github.com/nats-io/nats.go"
	"github.com/rivo/tview"
)

type SpelltextClient struct {
	Config      *config.Config
	Logger      *logging.Logger
	Nats        *nats.Conn
	App         *tview.Application
	PageManager *factory.PageManager
	User        *SpelltextUser
	Clients     *Clients
	Context 		*context.Context

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
	ChatClient  pbChat.ChatClient
	StoreClient pbStore.StoreClient
}

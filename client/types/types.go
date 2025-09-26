package types

import (
	"github.com/komadiina/spelltext/client/config"
	"github.com/komadiina/spelltext/client/factory"
	"github.com/komadiina/spelltext/client/registry"
	"github.com/komadiina/spelltext/utils/singleton/logging"
	"github.com/nats-io/nats.go"
	"github.com/rivo/tview"
)

type SpelltextClient struct {
	Config      *config.Config
	Servers     *registry.Registry
	Logger      *logging.Logger
	Nats        *nats.Conn
	App         *tview.Application
	PageManager *factory.PageManager
	User        *SpelltextUser

	NavigateTo func(pageKey string)
}

type SpelltextUser struct {
	Username string
}

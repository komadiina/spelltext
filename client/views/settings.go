package views

import (
	"github.com/komadiina/spelltext/client/constants"
	"github.com/komadiina/spelltext/client/types"
	"github.com/rivo/tview"
)

func AddSettingsPage(c *types.SpelltextClient) {
	c.PageManager.RegisterFactory(constants.PAGE_SETTINGS, func() tview.Primitive {
		form := tview.NewForm()

		return form
	}, nil, func() {})
}

package views

import (
	"github.com/komadiina/spelltext/client/types"
	"github.com/rivo/tview"
)

func AddCombatPage(c *types.SpelltextClient) {
	onClose := func() {}

	c.PageManager.RegisterFactory(COMBAT_PAGE, func() tview.Primitive {
		return tview.NewTextView()
	}, nil, onClose)
}

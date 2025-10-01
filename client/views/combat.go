package views

import (
	"github.com/komadiina/spelltext/client/types"
	"github.com/rivo/tview"
)

func AddCombatPage(c *types.SpelltextClient) {
	onClose := func() {}

	c.PageManager.RegisterFactory(COMBAT_PAGE, func() tview.Primitive {
		flex := tview.NewFlex().SetDirection(tview.FlexRow).SetFullScreen(true)
		flex.SetBorder(true).SetBorderPadding(1, 1, 5, 5).SetTitle(" combat ")
		return flex
	}, nil, onClose)
}

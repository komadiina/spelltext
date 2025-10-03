package views

import (
	types "github.com/komadiina/spelltext/client/types"
	"github.com/rivo/tview"
)

func AddInventoryPage(c *types.SpelltextClient) {
	onClose := func() {}

	c.PageManager.RegisterFactory(PAGE_INVENTORY, func() tview.Primitive {
		flex := tview.NewFlex().SetDirection(tview.FlexRow).SetFullScreen(true)
		flex.SetBorder(true).SetBorderPadding(1, 1, 5, 5).SetTitle(" inventory ")
		return flex
	}, nil, onClose)
}

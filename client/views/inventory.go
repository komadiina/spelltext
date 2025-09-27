package views

import (
	types "github.com/komadiina/spelltext/client/types"
	"github.com/rivo/tview"
)

func AddInventoryPage(c *types.SpelltextClient) {
	onClose := func() {}

	c.PageManager.RegisterFactory(INVENTORY_PAGE, func() tview.Primitive {
		return tview.NewTextView()
	}, nil, onClose)
}

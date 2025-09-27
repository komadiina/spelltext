package views

import (
	types "github.com/komadiina/spelltext/client/types"
	"github.com/rivo/tview"
)

func AddStorePage(c *types.SpelltextClient) {
	onClose := func() {}

	c.PageManager.RegisterFactory(STORE_PAGE, func() tview.Primitive {
		return tview.NewTextView()
	}, nil, onClose)
}

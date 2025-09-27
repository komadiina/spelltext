package views

import (
	types "github.com/komadiina/spelltext/client/types"
	"github.com/rivo/tview"
)

func AddProgressPage(c *types.SpelltextClient) {
	onClose := func() {}

	c.PageManager.RegisterFactory(PROGRESS_PAGE, func() tview.Primitive {
		return tview.NewTextView()
	}, nil, onClose)
}

package views

import (
	types "github.com/komadiina/spelltext/client/types"
	"github.com/rivo/tview"
)

func AddCharactersPage(c *types.SpelltextClient) {
	onClose := func() {}

	c.PageManager.RegisterFactory(CHARACTERS_PAGE, func() tview.Primitive {
		return tview.NewTextView()
	}, nil, onClose)
}

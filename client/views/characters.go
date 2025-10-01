package views

import (
	types "github.com/komadiina/spelltext/client/types"
	"github.com/rivo/tview"
)

func AddCharactersPage(c *types.SpelltextClient) {
	onClose := func() {}

	c.PageManager.RegisterFactory(CHARACTERS_PAGE, func() tview.Primitive {
		flex := tview.NewFlex().SetDirection(tview.FlexRow)
		flex.SetBorder(true).SetBorderPadding(1, 1, 5, 5).SetTitle(" armory ")
		flex = flex.SetFullScreen(true)
		return flex
	}, nil, onClose)
}

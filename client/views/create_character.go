package views

import (
	"github.com/rivo/tview"
	"github.com/komadiina/spelltext/client/types"
	"github.com/komadiina/spelltext/client/constants"
)

func AddCreateCharacterPage(c *types.SpelltextClient) {
	c.PageManager.RegisterFactory(constants.PAGE_CREATE_CHARACTER, func() tview.Primitive {
		return tview.NewFlex()
	}, nil, func() {})
}
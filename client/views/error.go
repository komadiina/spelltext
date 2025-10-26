package views

import (
	"github.com/komadiina/spelltext/client/constants"
	"github.com/komadiina/spelltext/client/types"
	"github.com/komadiina/spelltext/client/utils"
	"github.com/rivo/tview"
)

func AddErrorPage(c *types.SpelltextClient) {
	c.PageManager.RegisterFactory(constants.PAGE_ERROR, func() tview.Primitive {
		return utils.GenerateErrorPage(c, "")
	}, nil, func() {})
}

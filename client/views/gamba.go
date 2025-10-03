package views

import (
	types "github.com/komadiina/spelltext/client/types"
	"github.com/rivo/tview"
)

func AddGambaPage(c *types.SpelltextClient) {
	onClose := func() {}

	c.PageManager.RegisterFactory(PAGE_GAMBA, func() tview.Primitive {
		dummy := tview.NewTextView()
		dummy.SetText(PAGE_GAMBA)

		flex := tview.NewFlex().
			SetDirection(tview.FlexRow).
			SetFullScreen(true).
			AddItem(dummy, 0, 4, false)

		flex.SetBorder(true).SetBorderPadding(1, 1, 5, 5).SetTitle(" gamba ")
		return flex
	}, nil, onClose)
}

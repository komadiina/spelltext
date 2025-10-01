package views

import (
	"fmt"

	types "github.com/komadiina/spelltext/client/types"
	"github.com/rivo/tview"
)

func AddMainmenuPage(c *types.SpelltextClient) {
	c.PageManager.RegisterFactory(MAINMENU_PAGE, func() tview.Primitive {
		empty := tview.NewTextView()

		banner := tview.NewTextView().
			SetDynamicColors(true).
			SetText(fmt.Sprintf(`> welcome back, adventurer! [blue]%s[""] - isn't it?`, c.User.Username))

		navlist := tview.NewList().
			AddItem("characters", "", 'c', func() { c.PageManager.Push(CHARACTERS_PAGE, false) }).
			AddItem("inventory", "", 'i', func() { c.PageManager.Push(INVENTORY_PAGE, false) }).
			AddItem("progress", "", 'p', func() { c.PageManager.Push(PROGRESS_PAGE, false) }).
			AddItem("combat", "", 'b', func() { c.PageManager.Push(COMBAT_PAGE, false) }).
			AddItem("store", "", 's', func() { c.PageManager.Push(STORE_PAGE, false) }).
			AddItem("gamba", "", 'g', func() { c.PageManager.Push(GAMBA_PAGE, false) }).
			AddItem("chat", "", 'y', func() { c.PageManager.Push(CHAT_PAGE, false) }).
			AddItem("quit :(", "", 'q', func() { c.App.Stop() })

		flex := tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(banner, 1, 1, false).
			AddItem(empty, 1, 1, false).
			AddItem(navlist, 0, 2, true).
			SetFullScreen(true)

		flex.SetBorder(true).SetBorderPadding(1, 1, 5, 5).SetTitle(" menu ")

		return flex
	}, nil, func() { /* noop */ })
}

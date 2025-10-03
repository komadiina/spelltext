package views

import (
	"fmt"

	types "github.com/komadiina/spelltext/client/types"
	"github.com/rivo/tview"
)

func AddMainmenuPage(c *types.SpelltextClient) {
	c.PageManager.RegisterFactory(PAGE_MAINMENU, func() tview.Primitive {
		empty := tview.NewTextView()

		banner := tview.NewTextView().
			SetDynamicColors(true).
			SetText(fmt.Sprintf(`> welcome back, adventurer! [blue]%s[""] - isn't it?`, c.User.Username))

		navlist := tview.NewList().
			AddItem("characters", "", 'c', func() { c.PageManager.Push(PAGE_CHARACTERS, false) }).
			AddItem("inventory", "", 'i', func() { c.PageManager.Push(PAGE_INVENTORY, false) }).
			AddItem("progress", "", 'p', func() { c.PageManager.Push(PAGE_PROGRESS, false) }).
			AddItem("combat", "", 'b', func() { c.PageManager.Push(PAGE_COMBAT, false) }).
			AddItem("store", "", 's', func() { c.PageManager.Push(PAGE_STORE, false) }).
			AddItem("gamba", "", 'g', func() { c.PageManager.Push(PAGE_GAMBA, false) }).
			AddItem("chat", "", 'y', func() { c.PageManager.Push(PAGE_CHAT, false) }).
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

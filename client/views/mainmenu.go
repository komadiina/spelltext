package views

import (
	"fmt"

	"github.com/komadiina/spelltext/client/constants"
	types "github.com/komadiina/spelltext/client/types"
	"github.com/komadiina/spelltext/client/utils"
	"github.com/rivo/tview"
)

func AddMainmenuPage(c *types.SpelltextClient) {
	c.PageManager.RegisterFactory(constants.PAGE_MAINMENU, func() tview.Primitive {
		empty := tview.NewTextView()

		banner := tview.NewTextView().
			SetDynamicColors(true).
			SetText(fmt.Sprintf(`> welcome back, adventurer! [blue]%s[""] - isn't it?`, c.User.Username))

		navlist := tview.NewList().
			AddItem("characters", "", 'c', func() { c.PageManager.Push(constants.PAGE_CHARACTERS, false) }).
			AddItem("inventory", "", 'i', func() { c.PageManager.Push(constants.PAGE_INVENTORY, false) }).
			AddItem("progress", "", 'p', func() { c.PageManager.Push(constants.PAGE_PROGRESS, false) }).
			AddItem("combat", "", 'b', func() { c.PageManager.Push(constants.PAGE_COMBAT, false) }).
			AddItem("store", "", 's', func() { c.PageManager.Push(constants.PAGE_STORE, false) }).
			AddItem("gamba", "", 'g', func() { c.PageManager.Push(constants.PAGE_GAMBA, false) }).
			AddItem("chat", "", 'y', func() { c.PageManager.Push(constants.PAGE_CHAT, false) }).
			AddItem("quit :(", "", 'q', func() { c.App.Stop() })

		guide := tview.NewFlex().
			SetDirection(tview.FlexColumn)
		guide.SetBorder(true)

		characters, len := utils.AddNavGuide("c", "characters")
		guide.AddItem(characters, len, 1, false)

		inventory, len := utils.AddNavGuide("i", "inventory")
		guide.AddItem(inventory, len, 1, false)

		progress, len := utils.AddNavGuide("p", "progress")
		guide.AddItem(progress, len, 1, false)

		combat, len := utils.AddNavGuide("b", "combat")
		guide.AddItem(combat, len, 1, false)

		store, len := utils.AddNavGuide("s", "store")
		guide.AddItem(store, len, 1, false)

		gamba, len := utils.AddNavGuide("g", "gamba")
		guide.AddItem(gamba, len, 1, false)

		chat, len := utils.AddNavGuide("y", "chat")
		guide.AddItem(chat, len, 1, false)

		flex := tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(banner, 1, 1, false).
			AddItem(empty, 1, 1, false).
			AddItem(navlist, 0, 2, true).
			AddItem(guide, 3, 1, false).
			SetFullScreen(true)

		flex.SetBorder(true).SetBorderPadding(1, 1, 5, 5).SetTitle(" menu ")

		return flex
	}, nil, func() { /* noop */ })
}

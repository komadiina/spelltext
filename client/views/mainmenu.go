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
		banner := tview.NewTextView().
			SetDynamicColors(true).
			SetText(fmt.Sprintf(`> welcome back, adventurer! [blue]%s[""] - isn't it?`, c.User.Username))

		list := tview.NewList()
		list.
			SetHighlightFullLine(true).
			AddItem("- character", "preview your characters", 'a', func() { c.NavigateTo(constants.PAGE_CHARACTER) }).
			AddItem("- inventory", "peek at what's in your bags", 'i', func() { c.NavigateTo(constants.PAGE_INVENTORY) }).
			// AddItem("- progress", "see what you've accomplished", 'p', func() { c.NavigateTo(constants.PAGE_PROGRESS) }).
			AddItem("- combat", "THE proving grounds", 'c', func() { c.NavigateTo(constants.PAGE_COMBAT) }).
			AddItem("- store", "buy some mighty goods", 's', func() { c.NavigateTo(constants.PAGE_STORE) }).
			AddItem("- gamba", "gold gold gold\n\n\n", 'g', func() { c.NavigateTo(constants.PAGE_GAMBA) }).
			AddItem("- chat", "talk to some peeps", 'y', func() { c.NavigateTo(constants.PAGE_CHAT) }).
			AddItem("- quit", "done for today?", 'q', func() { c.App.Stop() })
		list.SetBorder(true).SetBorderPadding(1, 1, 5, 5)

		updates := tview.NewBox().SetTitle(" [::b]updates[::-] ").SetBorder(true).SetBorderPadding(1, 1, 5, 5)

		guide := utils.CreateGuide([]*types.UnusableHotkey{
			{Key: "a", Desc: "character"},
			{Key: "i", Desc: "inventory"},
			// {Key: "p", Desc: "progress"},
			{Key: "b", Desc: "combat"},
			{Key: "s", Desc: "store"},
			{Key: "g", Desc: "gamble"},
			{Key: "y", Desc: "chat"},
			{Key: "q", Desc: "quit"},
		})

		f2 := tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(list, 0, 1, true)

		tv := tview.NewTextView().
			SetDynamicColors(true).SetWrap(true).SetWordWrap(true)
		tv.SetBorder(true).SetBorderPadding(0, 0, 1, 1)

		if c.Storage.SelectedCharacter != nil {
			char := c.Storage.SelectedCharacter
			tv.SetText(fmt.Sprintf(`[selected character]%s> %s[""] - %s`, "\n", char.CharacterName, char.Hero.Name))
		} else {
			tv.SetText(`no character selected -- select one from the [blue]character[""] page`)
		}

		f2.AddItem(
			tv,
			4, 1, false,
		)

		flex := tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(banner, 1, 1, false).
			AddItem(tview.NewTextView(), 1, 1, false).
			AddItem(
				tview.NewFlex().SetDirection(tview.FlexColumn).
					AddItem(f2, 0, 1, true).
					AddItem(updates, 0, 1, false), 0, 1, true).
			AddItem(guide, 3, 1, false).
			SetFullScreen(true)

		flex.SetBorder(true).SetBorderPadding(1, 1, 5, 5).SetTitle(" [::b]menu[::-] ")

		return flex
	}, nil, func() { /* noop */ })
}

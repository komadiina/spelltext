package views

import (
	"github.com/komadiina/spelltext/client/constants"
	"github.com/komadiina/spelltext/client/functions"
	"github.com/komadiina/spelltext/client/types"
	"github.com/komadiina/spelltext/client/utils"
	"github.com/rivo/tview"
)

func AddCombatPage(c *types.SpelltextClient) {
	onClose := func() {}

	c.PageManager.RegisterFactory(constants.PAGE_COMBAT, func() tview.Primitive {
		flex := tview.NewFlex().SetDirection(tview.FlexRow).SetFullScreen(true)
		flex.SetBorder(true).SetBorderPadding(1, 1, 5, 5).SetTitle(" [::b]combat[::-] ")

		npcs := functions.GetAvailableNpcs(c)

		list := tview.NewList()
		list.
			SetBorder(true).
			SetBorderPadding(1, 1, 5, 5).
			SetTitle(`[#f1f1f1][::b] available npcs [::-][""]`)

		if len(npcs) == 0 {
			list.AddItem("woah...", "no one here", 'e', func() {
				c.PageManager.Pop()
				onClose()
			})
		} else {
			for _, npc := range npcs {
				list.AddItem("> "+utils.GetFullNpcName(npc), npc.NpcTemplate.Description, 0, func() {
					c.Logger.Info("npc selected", "npc", npc)
					c.Storage.Ministate.CurrentNpc = npc
					c.PageManager.Push(constants.PAGE_FIGHT, false)
				})
			}
		}

		details := tview.
			NewTextView().
			SetDynamicColors(true).
			SetTextAlign(tview.AlignLeft).
			SetWrap(true).
			SetWordWrap(true)

		details.
			SetBorder(true).
			SetBorderPadding(1, 1, 5, 5).
			SetTitle(`[#f1f1f1][::b] npc details [::-][""]`)

		list.SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
			npc := npcs[index]
			details.SetText(utils.PrintNpcDetails(npc))
		})

		flexTop := tview.NewFlex().SetDirection(tview.FlexColumn)
		flexTop.
			AddItem(list, 0, 1, true).
			AddItem(details, 0, 1, false)

		guide := utils.CreateGuide([]*types.UnusableHotkey{
			{Key: "↑/↓", Desc: "navigate"},
			{Key: "enter", Desc: "select"},
			{Key: "esc", Desc: "back"},
		})

		flex.
			AddItem(flexTop, 0, 1, true).
			AddItem(guide, 3, 1, false)

		return flex
	}, nil, onClose)
}

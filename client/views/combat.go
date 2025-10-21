package views

import (
	"fmt"

	"github.com/komadiina/spelltext/client/constants"
	"github.com/komadiina/spelltext/client/functions"
	"github.com/komadiina/spelltext/client/types"
	"github.com/komadiina/spelltext/client/utils"
	"github.com/rivo/tview"
)

func AddCombatPage(c *types.SpelltextClient) {
	onClose := func() {}

	c.PageManager.RegisterFactory(constants.PAGE_COMBAT, func() tview.Primitive {
		flex := tview.NewFlex().SetDirection(tview.FlexColumn).SetFullScreen(true)
		flex.SetBorder(true).SetBorderPadding(1, 1, 5, 5).SetTitle(" [::b]combat[::-] ")

		npcs := functions.GetAvailableNpcs(c)

		list := tview.NewList()
		list.SetBorder(true).SetTitle(`[#f1f1f1][::b] available npcs [::-][""]`)

		if len(npcs) == 0 {
			list.AddItem("woah...", "no one here", 'e', func() {
				c.PageManager.Pop()
				onClose()
			})
		} else {
			for _, npc := range npcs {
				list.AddItem("> "+utils.GetFullNpcName(npc), npc.NpcTemplate.Description, 0, func() {
					c.Logger.Info("npc selected", "npc", npc)
				})
			}
		}

		details := tview.NewTextView().SetDynamicColors(true).SetTextAlign(tview.AlignLeft)
		details.SetBorder(true).SetTitle(`[#f1f1f1][::b] npc details [::-][""]`)

		list.SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
			npc := npcs[index]
			c.Logger.Info("npc navigated to", "npc", npc)
			details.SetText(fmt.Sprintf("lv. %d %s", npc.Level, utils.GetFullNpcName(npc)))
		})

		flex.AddItem(list, 0, 1, true)
		flex.AddItem(details, 0, 1, true)

		return flex
	}, nil, onClose)
}

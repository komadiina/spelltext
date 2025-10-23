package views

import (
	"github.com/komadiina/spelltext/client/constants"
	"github.com/komadiina/spelltext/client/types"
	"github.com/komadiina/spelltext/client/utils"
	"github.com/rivo/tview"
)

func AddFightPage(c *types.SpelltextClient) {
	c.PageManager.RegisterFactory(constants.PAGE_FIGHT, func() tview.Primitive {
		flex := tview.NewFlex().SetDirection(tview.FlexRow).SetFullScreen(true)
		flex.SetBorder(true).SetBorderPadding(1, 1, 5, 5).SetTitle(" [::b]fight[::-] ")

		npc := c.Storage.Ministate.CurrentNpc
		c.Storage.Ministate.FightState = &types.NpcFightState{
			Npc: npc, CurrentHealth: int64(npc.NpcTemplate.HealthPoints) * int64(npc.HealthMultiplier),
		}

		npcDetails := tview.NewFlex().SetDirection(tview.FlexRow)
		npcDetails.SetBorder(true).SetBorderPadding(1, 1, 5, 5).SetTitle(" [::b]npc details[::-] ")
		nameTv := tview.NewTextView().SetText(utils.BoldText(utils.GetFullNpcName(npc))).SetDynamicColors(true)
		descTv := tview.NewTextView().SetText(npc.NpcTemplate.Description).SetDynamicColors(true)

		npcDetails.
			AddItem(nameTv, 1, 1, false).
			AddItem(descTv, 1, 1, false)

		flex.AddItem(npcDetails, 6, 1, false)

		return flex
	}, nil, func() {})
}

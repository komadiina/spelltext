package views

import (
	"fmt"

	"github.com/komadiina/spelltext/client/constants"
	"github.com/komadiina/spelltext/client/functions"
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

		npcHeader := tview.NewFlex().SetDirection(tview.FlexRow)
		npcHeader.SetBorder(true).SetBorderPadding(1, 1, 5, 5).SetTitle(" [::b]npc details[::-] ")
		nameTv := tview.NewTextView().SetText(utils.BoldText(utils.GetFullNpcName(npc))).SetDynamicColors(true)
		descTv := tview.NewTextView().SetText(npc.NpcTemplate.Description).SetDynamicColors(true)

		npcDetails := tview.NewFlex().SetDirection(tview.FlexRow)
		npcDetails.SetBorder(true).SetBorderPadding(1, 1, 2, 2).SetTitle(fmt.Sprintf(" [::b]%s[::-] ", utils.GetFullNpcName(npc)))

		const MAX_WIDTH uint8 = 32

		playerStatus := functions.InitEntityStatusFrame(32, *utils.GetDisplayStatsPlayer(c)...)
		npcStatus := functions.InitEntityStatusFrame(32, *utils.GetDisplayStatsNpc(npc)...)

		playerFrame := tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(playerStatus.FlHealth, 1, 0, false).
			AddItem(nil, 1, 0, false).
			AddItem(playerStatus.FlPower, 1, 0, false).
			AddItem(nil, 1, 0, false).
			AddItem(playerStatus.FlTextual, 1, 0, false)

		playerFrame.SetBorder(true).
			SetBorderPadding(1, 1, 2, 2).
			SetTitle(fmt.Sprintf(" [::b]%s[::-] ", c.Storage.SelectedCharacter.CharacterName))

		npcFrame := tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(npcStatus.FlHealth, 1, 0, false).
			AddItem(nil, 1, 0, false).
			AddItem(npcStatus.FlPower, 1, 0, false).
			AddItem(nil, 1, 0, false).
			AddItem(npcStatus.FlTextual, 1, 0, false)

		npcFrame.SetBorder(true).
			SetBorderPadding(1, 1, 2, 2).
			SetTitle(fmt.Sprintf(" [::b]%s[::-] ", utils.GetFullNpcName(npc)))

		flexStatusFrames := tview.NewFlex().
			SetDirection(tview.FlexColumn).
			AddItem(playerFrame, 0, 1, false).
			AddItem(nil, 0, 1, false).
			AddItem(npcFrame, 0, 1, false)

		npcHeader.
			AddItem(nameTv, 1, 1, false).
			AddItem(descTv, 1, 1, false)

		flex.
			AddItem(npcHeader, 6, 1, false).
			AddItem(flexStatusFrames, 9, 1, false)

		return flex
	}, nil, func() {})
}

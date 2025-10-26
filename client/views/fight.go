package views

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/komadiina/spelltext/client/constants"
	"github.com/komadiina/spelltext/client/functions"
	"github.com/komadiina/spelltext/client/types"
	"github.com/komadiina/spelltext/client/utils"
	clientUtils "github.com/komadiina/spelltext/client/utils"
	pbRepo "github.com/komadiina/spelltext/proto/repo"
	stUtils "github.com/komadiina/spelltext/utils"
	"github.com/rivo/tview"
)

func AddFightPage(c *types.SpelltextClient) {
	c.PageManager.RegisterFactory(constants.PAGE_FIGHT, func() tview.Primitive {
		rand.Seed(time.Now().UnixNano()) // deprecated but good enough

		var turnStarted bool = false
		var finished bool = false

		flex := tview.NewFlex().SetDirection(tview.FlexRow).SetFullScreen(true)
		flex.SetBorder(true).SetBorderPadding(1, 1, 5, 5).SetTitle(" [::b]fight[::-] ")

		npc := c.Storage.Ministate.CurrentNpc
		c.Storage.Ministate.FightState = &types.CbFightState{
			Npc:                 npc,
			NpcCurrentHealth:    int64(npc.NpcTemplate.HealthPoints) * int64(npc.HealthMultiplier),
			PlayerCurrentHealth: c.Storage.CharacterStats.HealthPoints,
			PlayerCurrentPower:  c.Storage.CharacterStats.PowerPoints,
		}

		// -------------- npc details
		npcHeader := tview.NewFlex().SetDirection(tview.FlexRow)
		npcHeader.SetBorder(true).SetBorderPadding(1, 1, 5, 5).SetTitle(" [::b]npc details[::-] ")
		nameTv := tview.NewTextView().SetText(clientUtils.BoldText(clientUtils.GetFullNpcName(npc))).SetDynamicColors(true)
		descTv := tview.NewTextView().SetText(npc.NpcTemplate.Description).SetDynamicColors(true)

		npcHeader.
			AddItem(nameTv, 1, 1, false).
			AddItem(descTv, 1, 1, false)

		npcDetails := tview.NewFlex().SetDirection(tview.FlexRow)
		npcDetails.SetBorder(true).SetBorderPadding(1, 1, 2, 2).SetTitle(fmt.Sprintf(" [::b]%s[::-] ", clientUtils.GetFullNpcName(npc)))

		// --------------- status frames
		playerStatus := functions.InitEntityStatusFrame(*clientUtils.GetDisplayStatsPlayer(c)...)
		npcStatus := functions.InitEntityStatusFrame(*clientUtils.GetDisplayStatsNpc(npc)...)

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
			SetTitle(fmt.Sprintf(" [::b]%s[::-] ", clientUtils.GetFullNpcName(npc)))

		flexStatusFrames := tview.NewFlex().
			SetDirection(tview.FlexColumn).
			AddItem(playerFrame, 0, 1, false).
			AddItem(nil, 0, 1, false).
			AddItem(npcFrame, 0, 1, false)

		// ----------------- combat log
		combatLog := tview.NewTextView().
			SetTextAlign(tview.AlignLeft).
			SetDynamicColors(true).
			ScrollToEnd()
		combatLog.
			SetBorder(true).
			SetBorderPadding(1, 1, 2, 2).
			SetTitle("[::b] combat log [::-]")

		functions.AddToCombatLog(combatLog, "combat started.")

		// ----------------- ability lists (player, npc)
		playerAbilities := stUtils.Map(
			functions.GetAbilities(c).Upgraded,
			func(pat *pbRepo.PlayerAbilityTree) *pbRepo.Ability {
				return pat.Ability
			})

		plAbList := tview.NewList()
		plAbList.SetBorder(true).SetBorderPadding(1, 1, 5, 5)

		for _, ab := range playerAbilities {
			plAbList.AddItem("> "+ab.GetName()+" ", "", 0, func() { // " " for a bit more visual clarity
				if turnStarted == false && finished == false {
					// todo: refactor all this to `functions` package
					if ab.PowerCost > uint64(c.Storage.Ministate.FightState.PlayerCurrentPower) {
						c.Logger.Warnf("unable to cast spell: not enough power. (cost=%d, have=%d)", ab.PowerCost, c.Storage.Ministate.FightState.PlayerCurrentPower)
						return
					}

					// lock turn
					turnStarted = true

					// calculate player damage
					dmgDone := functions.PlayerAttack(c, ab, c.Storage.Ministate.FightState)
					functions.RefreshStatusFrame(
						playerStatus,
						uint64(c.Storage.CharacterStats.HealthPoints), // keep
						uint64(c.Storage.CharacterStats.HealthPoints), // keep
						uint64(c.Storage.Ministate.FightState.PlayerCurrentPower),
						uint64(int(c.Storage.Ministate.FightState.PlayerCurrentPower)-int(ab.PowerCost)), // update power
					)

					functions.CombatLogTurn(
						combatLog,
						c.Storage.SelectedCharacter.CharacterName,
						ab.GetName(),
						utils.GetFullNpcName(c.Storage.Ministate.FightState.Npc),
						int(dmgDone),
					)

					functions.RefreshStatusFrame(
						npcStatus,
						uint64(functions.CalculateNpcStats(npc).HealthPoints),   // keep
						uint64(c.Storage.Ministate.FightState.NpcCurrentHealth), // was updated by functions.PlayerAttack, refresh
						0, // been attacked, no need to modify power
						0,
					)

					if c.Storage.Ministate.FightState.NpcCurrentHealth == 0 {
						functions.AddToCombatLog(combatLog, "you won! wow... congrats.")
						functions.AddToCombatLog(combatLog, fmt.Sprintf("you gain: [blue]%d xp[::-][white]", npc.NpcTemplate.BaseXpReward))
						finished = true
						functions.SubmitWin(c, npc)
						return
					}

					dmgDone = functions.NpcAttack(c, npc)
					functions.RefreshStatusFrame(
						playerStatus,
						uint64(c.Storage.CharacterStats.HealthPoints),
						uint64(c.Storage.Ministate.FightState.PlayerCurrentHealth),
						uint64(c.Storage.CharacterStats.PowerPoints),
						uint64(c.Storage.Ministate.FightState.PlayerCurrentPower),
					)
					functions.CombatLogTurn(
						combatLog,
						utils.GetFullNpcName(npc),
						"Attack",
						c.Storage.SelectedCharacter.CharacterName,
						int(dmgDone),
					)

					if c.Storage.Ministate.FightState.PlayerCurrentHealth == 0 {
						functions.AddToCombatLog(combatLog, "you lost. no surprise..")
						finished = true
						functions.SubmitLoss(c)
						return
					}

					// unlock turn
					turnStarted = false
				}
			})
		}

		plAbList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyCtrlF {
				functions.SubmitLoss(c)
				c.PageManager.Pop()
				return nil
			}

			return event
		})

		// ----------------- selected ability details
		abilityDetails := tview.NewTextView().
			SetWrap(true).
			SetWordWrap(true).
			SetDynamicColors(true)

		abilityDetails.
			SetBorder(true).
			SetBorderPadding(1, 1, 2, 2).
			SetTitle("[::b] ability details [::-]")

		plAbList.SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
			abilityDetails.SetText(functions.GetSpellDetails(playerAbilities[index]))
		})

		if len(playerAbilities) == 0 {
			return utils.GenerateErrorPage(c, "no abilities unlocked/available.")
		}

		abilityDetails.SetText(functions.GetSpellDetailsShort(playerAbilities[0]))

		// ------------------ guide
		guide := utils.CreateGuide([]*types.UnusableHotkey{
			{Key: "ctrl+f", Desc: "forfeit"},
			{Key: "enter", Desc: "play spell"},
			{Key: "↑/↓", Desc: "navigate"},
		}, true)

		flex.
			AddItem(npcHeader, 6, 1, false).
			AddItem(flexStatusFrames, 9, 1, false).
			AddItem(plAbList, 0, 1, true).
			AddItem(abilityDetails, 10, 1, false).
			AddItem(combatLog, 10, 1, false).
			AddItem(guide, 3, 1, false)

		return flex
	}, nil, func() {
		c.Storage.Ministate.FightState = nil
	})
}

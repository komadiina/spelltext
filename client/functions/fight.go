package functions

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/komadiina/spelltext/client/constants"
	"github.com/komadiina/spelltext/client/types"
	pbCombat "github.com/komadiina/spelltext/proto/combat"
	pbRepo "github.com/komadiina/spelltext/proto/repo"
	"github.com/rivo/tview"
)

func RedrawBar(current uint64, max uint64, maxBlocks uint8) string {
	pct := math.Round(float64(current) / float64(max) * float64(maxBlocks))
	numBlocks := int64(pct)

	// numBlocks := int64(float64(current) / float64(max))

	sb := strings.Builder{}
	for i := 0; i < int(numBlocks); i++ {
		sb.WriteString(constants.CHARACTER_HEALTH)
	}

	return sb.String()
}

func InitEntityStatusFrame(posStats ...uint64) *types.EntityStatusFrame {
	hp, pwr := posStats[0], posStats[1]

	// init health bar
	healthBar := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignRight)
	healthFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	healthTv := tview.NewTextView().
		SetDynamicColors(true).
		SetText(fmt.Sprintf(`[%s]HP: [::-]`, constants.TEXT_COLOR_HEALTH))
	hpTextLen := len("HP: ")
	hpBarText := RedrawBar(hp, hp, constants.MAX_STATUS_FRAME_WIDTH-uint8(hpTextLen))
	healthBar.SetText(fmt.Sprintf("[%s]%s[::-]", constants.TEXT_COLOR_HEALTH, hpBarText))
	healthFlex.AddItem(healthTv, hpTextLen, 1, false).
		AddItem(healthBar, 0, 1, false)

	// init power bar
	powerBar := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignRight)
	powerFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	powerTv := tview.NewTextView().
		SetDynamicColors(true).
		SetText(fmt.Sprintf(`[%s]PWR: [::-]`, constants.TEXT_COLOR_POWER))
	pwrTextLen := len("PWR: ")
	pwrBarText := RedrawBar(pwr, pwr, constants.MAX_STATUS_FRAME_WIDTH-uint8(pwrTextLen))
	powerBar.SetText(fmt.Sprintf("[%s]%s[::-]", constants.TEXT_COLOR_POWER, pwrBarText))
	powerFlex.AddItem(powerTv, pwrTextLen, 1, false).
		AddItem(powerBar, 0, 1, false)

	// init textual (HP | PWR) info
	textualFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	textualInfo := tview.NewTextView().SetTextAlign(tview.AlignCenter).SetDynamicColors(true)
	txtInfoText := fmt.Sprintf(`[::b][%s]%s[::-] | [%s]%s[::-][::-]`,
		constants.TEXT_COLOR_HEALTH, tview.Escape(fmt.Sprint("[", hp, "]")),
		constants.TEXT_COLOR_POWER, tview.Escape(fmt.Sprint("[", pwr, "]")))

	textualInfo.SetText(txtInfoText)
	textualFlex.AddItem(textualInfo, 0, 1, false)

	return &types.EntityStatusFrame{
		Health:      hp,
		Power:       pwr,
		BarHealth:   healthBar,
		BarPower:    powerBar,
		InfoTextual: textualInfo,
		FlHealth:    healthFlex,
		FlPower:     powerFlex,
		FlTextual:   textualFlex,
	}
}

func RefreshStatusFrame(esf *types.EntityStatusFrame, maxHealth uint64, newHp uint64, maxPower uint64, newPower uint64) {
	esf.Health = newHp
	esf.Power = uint64(newPower)

	hpTextLen := len("HP: ")
	hpBarText := RedrawBar(
		uint64(newHp),
		maxHealth,
		constants.MAX_STATUS_FRAME_WIDTH-uint8(hpTextLen),
	)

	pwrTextLen := len("PWR: ")
	pwrBarText := RedrawBar(
		newPower,
		maxPower,
		constants.MAX_STATUS_FRAME_WIDTH-uint8(pwrTextLen),
	)

	// update esf.BarHealth, esf.BarPower and esf.InfoTextual
	esf.BarHealth.SetText(fmt.Sprintf("[%s]%s[::-]", constants.TEXT_COLOR_HEALTH, hpBarText))
	esf.BarPower.SetText(fmt.Sprintf("[%s]%s[::-]", constants.TEXT_COLOR_POWER, pwrBarText))
	txtInfoText := fmt.Sprintf(`[::b][%s]%s[::-] | [%s]%s[::-][::-]`,
		constants.TEXT_COLOR_HEALTH, tview.Escape(fmt.Sprint("[", newHp, "]")),
		constants.TEXT_COLOR_POWER, tview.Escape(fmt.Sprint("[", newPower, "]")))
	esf.InfoTextual.SetText(txtInfoText)
}

func AddToCombatLog(cb *tview.TextView, text string) {
	fmt.Fprintf(cb, `[%s]%s[""][white][""]: %s%s`,
		constants.TEXT_COLOR_SPELL_UNLOCKED,
		tview.Escape(fmt.Sprint("[", time.Now().Format(time.TimeOnly), "]")),
		text,
		"\n",
	)

	cb.ScrollToEnd()
}

func CombatLogTurn(cb *tview.TextView, initiator string, spellName string, destEntity string, damage int) {
	str := fmt.Sprintf(`[%s]%s[""][white][""]: [%s]%s's[::-][white][""] [%s]%s[::-][white][""] hit [%s]%s[::-][white][""] for [%s]%d[::-][white][""] damage.%s`,
		constants.TEXT_COLOR_SPELL_UNLOCKED, tview.Escape(fmt.Sprint("[", time.Now().Format(time.TimeOnly), "]")),
		constants.TEXT_COLOR_NAME, initiator,
		constants.TEXT_COLOR_GOLD, spellName,
		constants.TEXT_COLOR_NAME, destEntity,
		constants.TEXT_COLOR_DAMAGE, damage, "\n",
	)

	fmt.Fprintf(cb, str)
	cb.ScrollToEnd()
}

func GetSpellDetailsShort(ability *pbRepo.Ability) string {
	return fmt.Sprintf(
		`[%s][::b]%s[::-][""][white][""]%s%s%s%sbase damage: [%s]%d[""][white][""]%spower cost: [%s]%d[""][white][""]`,
		constants.TEXT_COLOR_GOLD,
		ability.Name, "\n",
		ability.Description, "\n", "\n",
		constants.TEXT_COLOR_DAMAGE, ability.BaseDamage, "\n",
		constants.TEXT_COLOR_POWER, ability.PowerCost,
	)
}

func PlayerAttack(c *types.SpelltextClient, ab *pbRepo.Ability, fightState *types.CbFightState) uint64 {
	char := c.Storage.CharacterStats
	level := c.Storage.SelectedCharacter.Level

	variation := 0.06

	dmg := ab.BaseDamage

	spCoeff := uint64(ab.SpellpowerMultiplier)*uint64(char.Spellpower) + uint64(ab.SpMultPerlevel)*uint64(level)
	dmg += spCoeff

	strCoeff := uint64(ab.StrengthMultiplier)*uint64(char.Strength) + uint64(ab.StMultPerlevel)*uint64(level)
	dmg += strCoeff

	lower, upper := float64(dmg)-(float64(dmg)*variation), float64(dmg)+(float64(dmg)*variation)
	dmg = uint64(rand.Intn(int(upper)-int(lower)+1) + int(lower))

	if fightState.NpcCurrentHealth <= int64(dmg) {
		fightState.NpcCurrentHealth = 0
	} else {
		fightState.NpcCurrentHealth -= int64(dmg)
	}

	// update power
	if fightState.PlayerCurrentPower <= int64(ab.PowerCost) {
		fightState.PlayerCurrentPower = 0
	} else {
		fightState.PlayerCurrentPower -= int64(ab.PowerCost)
	}

	return dmg
}

func CalculateNpcDamage(npc *pbRepo.Npc) uint64 {
	variation := float32(npc.NpcTemplate.BaseDamage) * constants.NPC_DAMAGE_VARIATION_PCT
	dmg := npc.DamageMultiplier * float32(npc.NpcTemplate.BaseDamage)
	lower, upper := dmg-variation, dmg+variation

	return rand.Uint64()%uint64(lower) + uint64(upper)
}

func NpcAttack(c *types.SpelltextClient, npc *pbRepo.Npc) uint64 {
	dmg := CalculateNpcDamage(npc)

	if uint64(c.Storage.Ministate.FightState.PlayerCurrentHealth) < dmg { // overflow waiting to happen :))
		c.Storage.Ministate.FightState.PlayerCurrentHealth = 0
	} else {
		c.Storage.Ministate.FightState.PlayerCurrentHealth -= int64(dmg)
	}

	return dmg
}

func CalculateNpcStats(npc *pbRepo.Npc) *types.EntityStats {
	return &types.EntityStats{
		HealthPoints: int64(float32(npc.NpcTemplate.HealthPoints) * float32(npc.HealthMultiplier)),
		Damage:       int64(float32(npc.NpcTemplate.BaseDamage) * float32(npc.DamageMultiplier)),
	}
}

func SubmitLoss(c *types.SpelltextClient) {
	// TODO
}

func SubmitWin(c *types.SpelltextClient, npc *pbRepo.Npc) {
	resp, err := c.Clients.CombatClient.SubmitWin(*c.Context, &pbCombat.SubmitWinRequest{
		CharacterId: c.Storage.SelectedCharacter.CharacterId,
		NpcId:       npc.Id,
	})
	if err != nil {
		c.Logger.Error(err)
		return
	} else {
		c.Storage.SelectedCharacter = resp.NewCharacter
	}
}

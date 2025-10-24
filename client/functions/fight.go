package functions

import (
	"fmt"
	"math"
	"strings"

	"github.com/komadiina/spelltext/client/constants"
	"github.com/komadiina/spelltext/client/types"
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

func InitEntityStatusFrame(maxWidth uint8, posStats ...uint64) *types.EntityStatusFrame {
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
	hpBarText := RedrawBar(hp, hp, maxWidth-uint8(hpTextLen))
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
	pwrBarText := RedrawBar(pwr, pwr, maxWidth-uint8(pwrTextLen))
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
		Health:    hp,
		Power:     pwr,
		BarHealth: healthBar,
		BarPower:  powerBar,
		FlHealth:  healthFlex,
		FlPower:   powerFlex,
		FlTextual: textualFlex,
	}
}

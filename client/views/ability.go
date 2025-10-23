package views

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	pbRepo "github.com/komadiina/spelltext/proto/repo"

	"github.com/komadiina/spelltext/client/constants"
	"github.com/komadiina/spelltext/client/functions"
	"github.com/komadiina/spelltext/client/types"
	"github.com/komadiina/spelltext/client/utils"
	"github.com/rivo/tview"
)

func AddAbilityPage(c *types.SpelltextClient) {
	c.PageManager.RegisterFactory(constants.PAGE_ABILITY, func() tview.Primitive {
		flex := tview.NewFlex().SetDirection(tview.FlexRow).SetFullScreen(true)
		flex.SetBorder(true).SetBorderPadding(1, 1, 5, 5).SetTitle(" [::b]abilities[::-] ")

		tvDetails := tview.NewTextView().SetDynamicColors(true).SetWrap(true).SetWordWrap(true)
		tvDetails.SetBorder(true).SetBorderPadding(1, 1, 5, 5).SetTitle(" [::b]ability details[::-] ")
		available, locked, upgraded := functions.GetAbilities(c)
		all := append([]*pbRepo.Ability{}, append(*available, append(*locked, *upgraded...)...)...) // xdddd

		list := tview.NewList()
		list.SetBorder(true).SetBorderPadding(1, 1, 5, 5).SetTitle(" [::b]unlock abilities[::-] ")

		for _, ab := range *upgraded {
			list.AddItem(
				fmt.Sprint("> ", utils.ToColorTag(constants.TEXT_COLOR_SPELL_UNLOCKED), ab.Name, `[""][white][""]`),
				utils.PaintText(tcell.ColorGrey.String(), ab.Description), 0, func() {
					if err := functions.UpgradeAbility(c, ab); err != nil {
						c.Logger.Error(err)
					} else {

					}
				})
		}

		for _, ab := range *available {
			list.AddItem(
				fmt.Sprint("> ", utils.ToColorTag(constants.TEXT_COLOR_SPELL_AVAILABLE), ab.Name, `[""][white][""]`),
				utils.PaintText(tcell.ColorGrey.String(), ab.Description), 0, func() {
					if err := functions.UpgradeAbility(c, ab); err != nil {
						c.Logger.Error(err)
					} else {
						
					}
				})
		}

		for _, ab := range *locked {
			list.AddItem(
				fmt.Sprint("> ", utils.ToColorTag(constants.TEXT_COLOR_SPELL_LOCKED), ab.Name, `[""][white][""]`),
				utils.PaintText(tcell.ColorGrey.String(), ab.Description), 0, func() {
					if err := functions.UpgradeAbility(c, ab); err != nil {
						c.Logger.Error(err)
					} else {
						
					}
				})
		}

		list.SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
			tvDetails.SetText(functions.GetSpellDetails(all[index]))
		})

		// init
		tvDetails.SetText(functions.GetSpellDetails(all[0]))

		help := tview.NewTextView().
			SetText(functions.GetSpellDetailsHelp()).
			SetTextAlign(tview.AlignLeft).
			SetDynamicColors(true)

		help.SetBorder(true).SetBorderPadding(1, 1, 5, 5).SetTitle(" [::b]help[::-] ")

		flex.
			AddItem(list, 0, 1, true).
			AddItem(tvDetails, 0, 1, false).
			AddItem(help, 8, 1, false)

		return flex
	}, nil, func() {})
}

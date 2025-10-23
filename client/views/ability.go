package views

import (
	"fmt"

	pbRepo "github.com/komadiina/spelltext/proto/repo"
	generics "github.com/komadiina/spelltext/utils"

	"github.com/komadiina/spelltext/client/constants"
	"github.com/komadiina/spelltext/client/functions"
	"github.com/komadiina/spelltext/client/types"
	"github.com/komadiina/spelltext/client/utils"
	"github.com/rivo/tview"
)

func RenderList(c *types.SpelltextClient, ab []*pbRepo.Ability, color string, prefixFunc func(index int) string, callback func(*pbRepo.Ability, int), list *tview.List, clear bool) {
	if clear {
		list.Clear()
	}

	c.Logger.Info(ab)

	for index, ab := range ab {
		list.AddItem(
			fmt.Sprint("> ", prefixFunc(index), " ", utils.ToColorTag(color), ab.Name, `[""][white][""]`),
			ab.Description,
			0,
			func() { callback(ab, index) },
		)
	}
}

func UpdateAvailTpTv(tv *tview.TextView, c *types.SpelltextClient) *tview.TextView {
	tv.SetText(fmt.Sprintf("available talent points: %d", c.Storage.SelectedCharacter.UnspentPoints))
	return tv
}

func AddAbilityPage(c *types.SpelltextClient) {
	c.PageManager.RegisterFactory(constants.PAGE_ABILITY, func() tview.Primitive {
		flex := tview.NewFlex().SetDirection(tview.FlexRow).SetFullScreen(true)
		flex.SetBorder(true).SetBorderPadding(1, 1, 5, 5).SetTitle(" [::b]abilities[::-] ")

		tvDetails := tview.NewTextView().SetDynamicColors(true).SetWrap(true).SetWordWrap(true)
		tvDetails.SetBorder(true).SetBorderPadding(1, 1, 5, 5).SetTitle(" [::b]ability details[::-] ")

		availTp := tview.NewTextView().SetDynamicColors(true)
		availTp.SetBorder(true).SetBorderPadding(0, 0, 2, 2)
		UpdateAvailTpTv(availTp, c)

		lockedTp := tview.NewTextView().SetDynamicColors(true).SetWrap(true).SetWordWrap(true)
		lockedTp.SetBorder(true).SetBorderPadding(1, 1, 5, 5).SetTitle(" [::b]locked abilities[::-] ")

		list := tview.NewList()
		list.SetBorder(true).SetBorderPadding(1, 1, 5, 5).SetTitle(" [::b]unlock abilities[::-] ")

		tvInfo := tview.NewTextView().SetDynamicColors(true)

		resp, _ := functions.GetAbilities(c)
		available, locked, upgraded := resp.Available, resp.Locked, resp.Upgraded

		// _upgraded := []*pbRepo.Ability{}
		_upgraded := generics.Map(upgraded, func(a *pbRepo.PlayerAbilityTree) *pbRepo.Ability { return a.Ability })
		// for _, u := range upgraded {
		// 	_upgraded = append(_upgraded, u.Ability)
		// }

		all := append([]*pbRepo.Ability{}, append(available, append(locked, _upgraded...)...)...) // xdddd

		callback := func(a *pbRepo.Ability, _ int, newAbility bool) {
			if err := functions.UpgradeAbility(c, a, newAbility); err != nil {
				c.Logger.Error(err)
				tvInfo.SetText("[red]oops... an error occurred.[::-] try again later.")
			} else {
				UpdateAvailTpTv(availTp, c)
				if newAbility {
					c.PageManager.Pop()
					c.PageManager.Push(constants.PAGE_ABILITY, false)
				}
			}
		}

		RenderList(c, _upgraded, constants.TEXT_COLOR_SPELL_UNLOCKED,
			func(index int) string {
				a := (upgraded)[index]
				return fmt.Sprintf("(%d)", a.Level)
			},
			func(ab *pbRepo.Ability, index int) {
				callback(ab, index, false)
			}, list, true)

		RenderList(c, available, constants.TEXT_COLOR_SPELL_AVAILABLE,
			func(index int) string { return "0" },
			func(ab *pbRepo.Ability, index int) {
				callback(ab, index, true)
			}, list, false)

		RenderList(c, locked, constants.TEXT_COLOR_SPELL_LOCKED,
			func(index int) string { return "LCK" }, nil, list, false)

		list.SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
			tvDetails.SetText(functions.GetSpellDetails(all[index]))
		})

		tvDetails.SetText(functions.GetSpellDetails(all[0]))

		help := tview.NewTextView().
			SetText(functions.GetSpellDetailsHelp()).
			SetTextAlign(tview.AlignLeft).
			SetDynamicColors(true)

		help.SetBorder(true).SetBorderPadding(1, 1, 5, 5).SetTitle(" [::b]help[::-] ")

		flex.
			AddItem(tvInfo, 1, 0, false).
			AddItem(availTp, 3, 1, false).
			AddItem(list, 0, 1, true).
			AddItem(tvDetails, 0, 1, false).
			AddItem(help, 8, 1, false)

		return flex
	}, nil, func() {})
}

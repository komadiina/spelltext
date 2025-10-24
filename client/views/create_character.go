package views

import (
	"fmt"

	"github.com/komadiina/spelltext/client/constants"
	"github.com/komadiina/spelltext/client/functions"
	"github.com/komadiina/spelltext/client/types"
	pbRepo "github.com/komadiina/spelltext/proto/repo"
	generics "github.com/komadiina/spelltext/utils"
	"github.com/rivo/tview"
)

func AddCreateCharacterPage(c *types.SpelltextClient) {
	c.PageManager.RegisterFactory(constants.PAGE_CREATE_CHARACTER, func() tview.Primitive {
		flex := tview.NewFlex().SetDirection(tview.FlexColumn).SetFullScreen(true)
		flex.SetBorder(true).SetBorderPadding(2, 2, 5, 5).SetTitle(" [::b]character creation[::-] ")

		form := tview.NewForm()
		form.SetBorder(true).SetBorderPadding(2, 2, 5, 5).SetTitle(" [::b]form[::-] ")

		heroDetails := tview.NewTextView().SetDynamicColors(true)
		heroDetails.SetBorder(true).SetBorderPadding(2, 2, 5, 5).SetTitle(" [::b]hero details[::-] ")

		var newCharacter *pbRepo.Character = &pbRepo.Character{}

		heroes := functions.ListHeroes(c)
		listHeroes := tview.NewList()
		for _, hero := range heroes {
			listHeroes.AddItem(fmt.Sprintf("> %s", hero.Name), "", 0, func() {
				newCharacter.Hero = hero
			})
		}

		form.
			AddInputField("character name: ", "", 20, nil, func(text string) {
				newCharacter.CharacterName = text
			}).
			AddDropDown(
				"hero",
				generics.Map(heroes, func(h *pbRepo.Hero) string { return h.Name }),
				0,
				func(option string, optionIndex int) {
					c.Logger.Debugf("%d %+v", optionIndex, heroes[optionIndex])
					newCharacter.Hero = heroes[optionIndex]
					heroDetails.SetText(newCharacter.Hero.Name)
				}).
			AddButton("create", func() {
				functions.CreateCharacter(newCharacter, c)
			})

		return form
	}, nil, func() {})
}

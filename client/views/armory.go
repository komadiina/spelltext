package views

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/komadiina/spelltext/client/constants"
	"github.com/komadiina/spelltext/client/functions"
	types "github.com/komadiina/spelltext/client/types"
	"github.com/komadiina/spelltext/client/utils"
	pb "github.com/komadiina/spelltext/proto/armory"
	"github.com/rivo/tview"
)

type CharacterDetailsView struct {
	Name     *tview.TextView
	Level    *tview.TextView
	Class    *tview.TextView
	Currency *tview.TextView
}

func (d *CharacterDetailsView) Update(c *pb.TCharacter) {
	d.Name.SetText(fmt.Sprintf("name: %s", c.Name))
	d.Level.SetText(fmt.Sprintf(`level: [blue]%d[""]`, c.Level))
	d.Class.SetText(fmt.Sprintf(`class: %s`, c.Class))
	d.Currency.SetText(fmt.Sprintf(`[yellow]%dg[""] | [orange]%dt[""]`, c.Gold, c.Tokens))
}

func AddArmoryPage(c *types.SpelltextClient) {
	onClose := func() {}

	c.PageManager.RegisterFactory(constants.PAGE_ARMORY, func() tview.Primitive {
		flex := tview.NewFlex().SetDirection(tview.FlexRow)
		flex.SetBorder(true).SetBorderPadding(1, 1, 5, 5).SetTitle(" armory ")
		flex = flex.SetFullScreen(true)

		header := tview.NewTextView().SetText("select one from available characters: ")
		characters := tview.NewList()

		var uid uint64 = 1 // TODO
		chars, err := functions.GetCharacters(uid, c)
		if err != nil {
			c.Logger.Error(err)
			chars.Characters = make([]*pb.TCharacter, 0)
		}

		stored := []*pb.TCharacter{}
		panel := tview.NewFlex().SetDirection(tview.FlexRow)
		details := CharacterDetailsView{
			tview.NewTextView().SetDynamicColors(true),
			tview.NewTextView().SetDynamicColors(true),
			tview.NewTextView().SetDynamicColors(true),
			tview.NewTextView().SetDynamicColors(true),
		}

		panel.SetBorder(true).SetBorderPadding(0, 1, 1, 1).SetTitle(" character details ").SetTitleAlign(tview.AlignLeft)
		panel = panel.
			AddItem(details.Name, 1, 1, false).
			AddItem(details.Level, 1, 1, false).
			AddItem(details.Class, 1, 1, false).
			AddItem(details.Currency, 1, 1, false)

		charSelected, ok := c.AppStorage[constants.SELECTED_CHARACTER].(*pb.TCharacter)
		if !ok {
			charSelected = &pb.TCharacter{Name: "none"}
		}

		selected := tview.NewTextView().
			SetText(
				fmt.Sprintf(`+++ currently selected: [orange]%s[""]!`, charSelected.Name)).
			SetDynamicColors(true)

		for _, character := range chars.Characters {
			stored = append(stored, character)
			characters.AddItem("-> "+character.Name, "", 0, func() {
				c.AppStorage[constants.SELECTED_CHARACTER] = character

				mod := tview.NewModal().SetText(fmt.Sprintf("character: %s", character.Name)).AddButtons([]string{"select", "delete", "cancel"})
				mod.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
					switch buttonIndex {
					case 0: // select
						c.AppStorage[constants.SELECTED_CHARACTER] = character
						c.App.SetRoot(c.PageManager.Pages, true).EnableMouse(true)
						selected.SetText(fmt.Sprintf(`+++ currently selected: [orange]%s[""] (lv. %d %s)`, character.Name, character.Level, character.Class))
						return
					case 1: // delete
						functions.DeleteCharacter(character, c)
						c.App.SetRoot(c.PageManager.Pages, true).EnableMouse(true)
						c.PageManager.Pop()
						c.NavigateTo(constants.PAGE_ARMORY)
						return
					case 2:
						c.App.SetRoot(c.PageManager.Pages, true).EnableMouse(true)
						return
					}
				})

				mod.SetBackgroundColor(tcell.ColorDarkSlateGrey)
				c.App.SetRoot(mod, true).EnableMouse(true)
			})
		}

		characters.SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
			details.Update(stored[index])
		})

		characters.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyCtrlA {
				c.Logger.Info("create character command pressed.")
			}

			return event
		})

		characters.SetBorder(true).SetTitle(" characters ").SetTitleAlign(tview.AlignLeft).SetBorderPadding(1, 1, 2, 2)

		guide := tview.NewFlex().
			SetDirection(tview.FlexColumn).
			AddItem(tview.NewTextView().SetText(" keymap legend: "), 0, 1, false)
		guide.SetBorder(true)

		back, len := utils.AddNavGuide("esc", "back")
		guide.AddItem(back, len, 1, false)

		add, len := utils.AddNavGuide("ctrl+a", "create new character")
		guide.AddItem(add, len, 1, false)

		enter, len := utils.AddNavGuide("enter", "character menu")
		guide.AddItem(enter, len, 1, false)

		flex.
			AddItem(panel, 6, 1, false).
			AddItem(tview.NewTextView(), 1, 1, false).
			AddItem(selected, 1, 1, false).
			AddItem(header, 1, 1, false).
			AddItem(tview.NewTextView(), 1, 1, false).
			AddItem(characters, 0, 1, true).
			AddItem(tview.NewTextView(), 0, 3, false).
			AddItem(guide, 3, 1, false)

		return flex
	}, nil, onClose)
}

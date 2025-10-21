package utils

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/komadiina/spelltext/client/types"
	pbRepo "github.com/komadiina/spelltext/proto/repo"
	"github.com/rivo/tview"
)

func AddNavGuide(shortcut string, name string) (*tview.TextView, int) {
	str := fmt.Sprintf(" [%s] %s ", strings.ToUpper(shortcut), name)
	tv := tview.NewTextView().SetText(str)
	return tv, len(str)
}

func CreateModal(title string, message string, c *types.SpelltextClient, onClose func()) *tview.Modal {
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"OK"})

	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if onClose != nil {
			onClose()
		} else {
			c.PageManager.Pop()
			c.App.SetRoot(c.PageManager.Pages, true).EnableMouse(true)
		}
	})

	modal.SetBackgroundColor(tcell.ColorDarkSlateGrey)
	modal.SetTitle(title)

	return modal
}

func UpdateGold(tv *tview.TextView, format string, delta int64, c *types.SpelltextClient) *tview.TextView {
	char := c.Storage.SelectedCharacter

	char.Gold = uint64(int64(char.Gold) + delta)
	c.Storage.SelectedCharacter = char

	return tv.SetText(fmt.Sprintf(format, char.Gold))
}

func UpdateCharacter(old *pbRepo.Character, new *pbRepo.Character, c *types.SpelltextClient) {
	c.Storage.SelectedCharacter = new
}

func UpdateCharacterFunc(char *pbRepo.Character, c *types.SpelltextClient, f func(*pbRepo.Character) *pbRepo.Character) {
	c.Storage.SelectedCharacter = f(char)
}

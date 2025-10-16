package utils

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/komadiina/spelltext/client/constants"
	"github.com/komadiina/spelltext/client/types"
	pbArmory "github.com/komadiina/spelltext/proto/armory"
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
	char := c.AppStorage[constants.SELECTED_CHARACTER].(*pbArmory.TCharacter)

	char.Gold = uint64(int64(char.Gold) + delta)
	c.AppStorage[constants.SELECTED_CHARACTER] = char

	return tv.SetText(fmt.Sprintf(format, char.Gold))
}

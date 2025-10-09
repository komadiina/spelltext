package utils

import (
	"fmt"
	"strings"

	"github.com/rivo/tview"
)

func AddNavGuide(shortcut string, name string) (tview.Primitive, int) {
	str := fmt.Sprintf(" [%s] %s ", strings.ToUpper(shortcut), name)
	tv := tview.NewTextView().SetText(str)
	return tv, len(str)
}

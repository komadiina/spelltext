package views

import "github.com/rivo/tview"

func STNewFlex() *tview.Flex {
	flex := tview.NewFlex().SetFullScreen(true)
	flex.SetBorderPadding(5, 5, 10, 10)
	return flex
}

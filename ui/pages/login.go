package pages

import "github.com/rivo/tview"

func GenerateLoginPage(app *tview.Application, pages *tview.Pages) *tview.Pages {
	username := tview.NewInputField()

	loginButton := tview.NewButton("login").SetSelectedFunc(func() {
		pages.SwitchToPage("chat")
	})

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(username, 0, 1, true).
		AddItem(loginButton, 0, 1, false)
	
	flex.SetBorder(true).SetTitle(" login ")

	pages.AddPage("login", flex, true, true)
	return pages
}

package views

import (
	"strings"

	"github.com/komadiina/spelltext/client/functions"
	types "github.com/komadiina/spelltext/client/types"
	"github.com/rivo/tview"
)

func AddLoginPage(c *types.SpelltextClient) {
	c.PageManager.RegisterFactory(PAGE_LOGIN, func() tview.Primitive {
		header1 := "> spelltext v0.2.0"
		header2 := "> https://github.com/komadiina/spelltext"

		form := tview.NewForm().
			AddTextView("", "[yellow]"+header1, 0, 1, true, false).
			AddTextView("", "[blue]"+header2, 0, 1, true, false).
			AddTextView("", "", 0, 1, true, false).
			AddInputField("Username: ", "", len(header2), nil, func(text string) {
				*c.User = functions.GetUserByUsername(strings.ToLower(text))
			}).
			AddPasswordField("Password: ", "", len(header2), '*', func(text string) {}).
			AddButton("Login", func() {
				c.PageManager.Push(PAGE_MAINMENU, false)
			})

		return form
	}, nil, func() { /* noop */ })
}

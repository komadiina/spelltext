package views

import (
	"strings"

	"github.com/komadiina/spelltext/client/constants"
	"github.com/komadiina/spelltext/client/functions"
	types "github.com/komadiina/spelltext/client/types"
	"github.com/rivo/tview"
)

func AddLoginPage(c *types.SpelltextClient) {
	c.PageManager.RegisterFactory(constants.PAGE_LOGIN, func() tview.Primitive {
		header1 := "> spelltext v0.3.0"
		header2 := "> https://github.com/komadiina/spelltext"

		username := tview.NewInputField().
			SetLabel("username: ").
			SetFieldWidth(len(header2)).
			SetChangedFunc(func(text string) { *c.User = functions.GetUserByUsername(strings.ToLower(text)) })

		var password string

		form := tview.NewForm().
			AddTextView("", "[yellow]"+header1, 0, 1, true, false).
			AddTextView("", "[blue]"+header2, 0, 1, true, false).
			AddTextView("", "", 0, 1, true, false).
			AddFormItem(username).
			AddPasswordField("password: ", "", len(header2), '*', func(text string) {
				password = strings.Trim(text, " ")
			}).
			AddButton("login", func() {
				uname := strings.ToLower(username.GetText())
				functions.LoginUser(c, uname, password)
				c.Storage.Ministate.Username = uname
				functions.SetLastSelectedCharacter(c)
				c.PageManager.Push(constants.PAGE_MAINMENU, false)
			})

		flex := STNewFlex().AddItem(form, 0, 1, true)
		flex.SetBorder(true).SetTitle(" [::b]login[::-] ")

		return flex
	}, nil, func() { /* noop */ })
}

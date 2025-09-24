package main

import (
	pg "github.com/komadiina/spelltext/ui/pages"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()
	pages := tview.NewPages()

	pages = pg.GenerateChat(app, pages)
	pages = pg.GenerateLoginPage(app, pages)

	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

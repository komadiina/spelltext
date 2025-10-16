package views

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/komadiina/spelltext/client/constants"
	"github.com/komadiina/spelltext/client/functions"
	types "github.com/komadiina/spelltext/client/types"
	pbArmory "github.com/komadiina/spelltext/proto/armory"
	"github.com/rivo/tview"
)

func AddInventoryPage(c *types.SpelltextClient) {
	onClose := func() {}

	c.PageManager.RegisterFactory(constants.PAGE_INVENTORY, func() tview.Primitive {
		flex := tview.NewFlex().SetDirection(tview.FlexRow).SetFullScreen(true)
		flex.SetBorder(true).SetBorderPadding(1, 1, 5, 5).SetTitle(" inventory ")

		char := c.AppStorage[constants.SELECTED_CHARACTER]
		if char == nil {
			flex.AddItem(tview.NewTextView().SetText("no character selected"), 0, 1, false)
			return flex
		} else {
			char := char.(*pbArmory.TCharacter)
			tv := tview.NewTextView().SetText(fmt.Sprintf("browsing %s's inventory", char.Name))
			tv.SetBackgroundColor(tcell.ColorSlateGrey).SetBorderPadding(1, 1, 2, 2)
			flex.AddItem(tv, 3, 1, false).AddItem(tview.NewTextView().SetWrap(true).SetWordWrap(true), 1, 1, false)
		}

		items := functions.GetBackpackItems(c).GetItems()

		if len(items) == 0 {
			flex.AddItem(tview.NewTextView().SetText("no items in backpack").SetWrap(true).SetWordWrap(true), 0, 1, false)
		} else {
			table := tview.NewTable().SetSeparator('|')
			table.SetBorder(true)
			table = functions.MakeVendorTableHeader(table)

			for idx, item := range items {
				// tv := tview.NewTextView().SetWrap(true).SetWordWrap(true).SetDynamicColors(true)
				// name := item.Prefix + " " + item.ItemTemplate.Name + " " + item.Suffix
				// tv.SetBorder(true).SetBorderPadding(0, 0, 2, 2).SetTitle(" " + strings.Trim(name, " ") + " ")
				// tv.SetText(
				// 	fmt.Sprintf(
				// 		"HP:\t\t%d\nPWR:\t%d\nSTR:\t%d\nSP:\t\t%d\nDMG:\t%d\nARMOR:\t%d\n%s",
				// 		item.Health, item.Power, item.Strength, item.Spellpower, item.BonusDamage, item.BonusArmor, item.ItemTemplate.Description,
				// 	),
				// )

				// flex.AddItem(tv, 9, 1, false)
				table.
					Select(1, 0).
					SetFixed(1, 0).
					SetDoneFunc(func(key tcell.Key) {
						switch key {
						case tcell.KeyEscape:
							c.PageManager.Pop() // todo: move focus to flex
						case tcell.KeyEnter:
							table.SetSelectable(true, false)
							// TODO: display item detailed stats (in modal)
						}
					}).
					SetSelectedFunc(func(row, column int) {
						table.SetSelectable(true, false)
					})
				table.SetEvaluateAllRows(true)
				table.SetBorderPadding(1, 1, 5, 5)

				table = functions.MakeInventoryTableRow(idx+1, item, c, table)
			}

			flex.AddItem(table, 0, 1, true)
		}

		return flex
	}, nil, onClose)
}

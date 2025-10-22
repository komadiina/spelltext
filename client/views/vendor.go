package views

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/komadiina/spelltext/client/constants"
	"github.com/komadiina/spelltext/client/functions"
	types "github.com/komadiina/spelltext/client/types"
	"github.com/komadiina/spelltext/client/utils"
	pbStore "github.com/komadiina/spelltext/proto/store"
	"github.com/rivo/tview"
)

func UpdateBasket(basket *[]*pbStore.Item, tv *tview.TextView, c *types.SpelltextClient) {
	var totalGold uint64 = 0
	for _, item := range *basket {
		totalGold += item.GetGoldPrice()
	}

	color := "yellow"
	if totalGold > c.Storage.SelectedCharacter.Gold {
		color = "red"
	}

	tv.SetText(fmt.Sprintf(`basket price: [%s]%dg[""]`, color, totalGold))
}

func AddVendorPage(c *types.SpelltextClient) {
	onClose := func() {}

	c.PageManager.RegisterFactory(constants.PAGE_VENDOR, func() tview.Primitive {
		if c.Storage.SelectedCharacter == nil {
			f := tview.NewFlex().SetFullScreen(true)
			f.SetBorder(true).SetBorderPadding(1, 1, 5, 5).SetTitle(" [::b]hello?[::-] ")

			return f.AddItem(tview.NewTextView().
				SetText("no character selected. select a character from the character page, and come back... dummy"), 0, 1, false)
		}

		basket := make([]*pbStore.Item, 0)
		totals := tview.NewFlex().SetDirection(tview.FlexRow)
		basketPrice := tview.NewTextView().SetDynamicColors(true).SetText(`basket price: [yellow]0g[""]`)
		basketPrice.SetBorder(true).SetBorderPadding(1, 1, 2, 2)

		availableGold := tview.
			NewTextView().
			SetDynamicColors(true).
			SetText(fmt.Sprintf(`available gold: [yellow]%d[""]`, c.Storage.SelectedCharacter.Gold))
		availableGold.SetBorder(true).SetBorderPadding(1, 1, 2, 2)

		totals.AddItem(basketPrice, 5, 1, false).AddItem(availableGold, 5, 1, false)

		vendor := tview.NewTextView().
			SetDynamicColors(true).
			SetText(fmt.Sprintf(`[blue]%v[""]'s wares`, c.Storage.SelectedVendor.Name))

		flex := STNewFlex().AddItem(vendor, 1, 1, false).SetDirection(tview.FlexRow)
		flex.SetBorder(true).SetTitle(fmt.Sprintf(" [::b]vendor - %s[::-] ", c.Storage.SelectedVendor.Name))

		flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyCtrlB {
				char := c.Storage.SelectedCharacter
				if err := functions.BuyItems(basket, char, c); err != nil {
					c.Logger.Error(err)
				}

				// reset basket, reload gold available
				go func() {
					c.App.QueueUpdateDraw(func() {
						availableGold.SetText(fmt.Sprintf(`available gold: [yellow]%d[""]`, char.GetGold()))
					})
				}()

				char = c.Storage.SelectedCharacter
				m := utils.CreateModal(
					"purchase successful",
					fmt.Sprintf("you've bought %d items (remaining gold: %d)", len(basket), char.GetGold()),
					c,
					nil,
				)
				basket = make([]*pbStore.Item, 0)

				c.App.SetRoot(m, true).EnableMouse(true)
			}

			return event
		})

		table := tview.NewTable().SetSeparator('|')
		table.SetBorder(true)

		go func() {
			c.App.QueueUpdateDraw(func() {
				functions.MakeVendorTableHeader(table)
			})
		}()

		list := tview.NewList()
		bought := make(map[uint64]bool)

		resp, err := c.Clients.StoreClient.ListVendorItems(
			*c.Context,
			&pbStore.StoreListVendorItemRequest{
				VendorId: c.Storage.SelectedVendor.Id,
			},
		)

		if err != nil {
			c.Logger.Error(err)
		}

		if len(resp.Items) == 0 {
			list.AddItem("woah...", "stock is empty. at the moment. check back later.", 'e', func() {})
		} else {
			for idx, item := range resp.Items {
				table = functions.MakeVendorTableRow(idx+1, item, table)
			}
		}

		details := tview.NewTextView().
			SetText(fmt.Sprintf(`total items: %v`, len(resp.Items)))
		flex.AddItem(details, 1, 0, false)

		table.
			Select(1, 0).
			SetFixed(1, 0).
			SetDoneFunc(func(key tcell.Key) {
				if key == tcell.KeyEscape {
					c.PageManager.Pop()
				}
				if key == tcell.KeyEnter {
					table.SetSelectable(true, false)
				}
			}).
			SetSelectedFunc(func(row, column int) {
				if bought[resp.Items[row-1].GetId()] {
					// remove from basket
					for idx, item := range basket {
						if item.GetId() == resp.Items[row-1].GetId() {
							basket = append(basket[:idx], basket[idx+1:]...)
							break
						}
					}

					// unmark red
					table.GetCell(row, 0).SetTextColor(constants.COLOR_NAME)
					UpdateBasket(&basket, basketPrice, c)
					bought[resp.Items[row-1].GetId()] = false
					return
				} else {
					// add to basket
					bought[resp.Items[row-1].GetId()] = true
					basket = append(basket, resp.Items[row-1])
					UpdateBasket(&basket, basketPrice, c)
					table.GetCell(row, 0).SetTextColor(constants.COLOR_GOLD)
					table.SetSelectable(true, false)
				}
			})

		table.SetEvaluateAllRows(true)
		table.SetBorderPadding(1, 1, 5, 5)

		guide := utils.CreateGuide([]*types.UnusableHotkey{
			{Key: "enter", Desc: "add to basket"},
			{Key: "ctrl+b", Desc: "buy"},
			{Key: "esc", Desc: "back"},
		})

		flex = flex.
			AddItem(table, 0, 1, true).
			AddItem(totals, 0, 1, false).
			AddItem(guide, 3, 1, false).
			SetFullScreen(true)

		return flex
	}, nil, onClose)
}

package views

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/komadiina/spelltext/client/constants"
	"github.com/komadiina/spelltext/client/functions"
	types "github.com/komadiina/spelltext/client/types"
	pb "github.com/komadiina/spelltext/proto/store"
	"github.com/rivo/tview"
)

func AddVendorPage(c *types.SpelltextClient) {
	onClose := func() {}

	c.PageManager.RegisterFactory(constants.PAGE_VENDOR, func() tview.Primitive {
		c.Logger.Info("loading vendor page...")

		vendor := tview.NewTextView().
			SetDynamicColors(true).
			SetText(fmt.Sprintf(`[blue]%v[""]'s wares`, c.AppStorage[SELECTED_VENDOR_NAME]))

		flex := STNewFlex().AddItem(vendor, 1, 1, false).SetDirection(tview.FlexRow)

		table := tview.NewTable().SetSeparator('|')
		table.SetBorder(true)
		table = functions.MakeVendorTableHeader(table)

		list := tview.NewList()

		resp, err := c.Clients.StoreClient.ListVendorItems(
			*c.Context,
			&pb.StoreListVendorItemRequest{
				VendorId: c.AppStorage[SELECTED_VENDOR_ID].(uint64),
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

		details := tview.NewTextView().SetText(fmt.Sprintf(`total items: %v`, len(resp.Items)))
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
				table.GetCell(row, 0).SetTextColor(tcell.ColorRed)
				table.SetSelectable(false, false)
			}).
			SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
				c.Logger.Info("key captured", "key", event.Key(), "rune", string(event.Rune()))
				return event
			})

		table.SetEvaluateAllRows(true)
		table.SetBorderPadding(1, 1, 5, 5)

		flex = flex.AddItem(table, 0, 1, true)

		return flex
	}, nil, onClose)
}

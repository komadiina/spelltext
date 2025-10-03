package views

import (
	types "github.com/komadiina/spelltext/client/types"
	pb "github.com/komadiina/spelltext/proto/store"
	"github.com/rivo/tview"
)

func AddVendorPage(c *types.SpelltextClient) {
	onClose := func() {}

	c.PageManager.RegisterFactory(PAGE_VENDOR, func() tview.Primitive {
		var vendor pb.Vendor

		flex := tview.NewFlex()
		list := tview.NewList()

		resp, err := c.Clients.StoreClient.ListVendorItems(
			*c.Context,
			&pb.StoreListVendorItemRequest{
				VendorId: c.AppStorage[SELECTED_VENDOR_ID].(uint64),
			},
		)

		if err != nil {
			list.AddItem("no items fetched :(", "check storeserver logs", 'e', func() { onClose() })
		} else {
			for idx, item := range resp.Items {
				list.AddItem(item.GetName(), item.GetDescription(), rune(idx), func() {})
			}
		}

		flex = flex.
			SetDirection(tview.FlexRow).
			AddItem(list, 0, 1, true).
			SetFullScreen(true)

		flex.SetBorder(true).SetBorderPadding(1, 1, 5, 5).SetTitle(vendor.VendorName)

		return flex
	}, nil, onClose)
}

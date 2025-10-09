package views

import (
	"github.com/komadiina/spelltext/client/constants"
	types "github.com/komadiina/spelltext/client/types"
	pb "github.com/komadiina/spelltext/proto/store"
	"github.com/rivo/tview"
)

func AddStorePage(c *types.SpelltextClient) {
	onClose := func() {}

	c.PageManager.RegisterFactory(constants.PAGE_STORE, func() tview.Primitive {
		flex := tview.NewFlex()
		list := tview.NewList()

		resp, err := c.Clients.StoreClient.ListVendors(*c.Context, &pb.StoreListVendorRequest{})
		if err != nil {
			c.Logger.Error(err)
		}

		if len(resp.Vendors) == 0 {
			list.AddItem("woah...", "stock is empty. at the moment.", 'e', func() { onClose() })
		} else {
			for idx, vendor := range resp.Vendors {
				list.AddItem(vendor.GetVendorName(), vendor.GetVendorWareDescription(), rune(idx), func() {
					c.AppStorage[SELECTED_VENDOR_ID] = vendor.GetVendorId()
					c.AppStorage[SELECTED_VENDOR_NAME] = vendor.GetVendorName()
					c.NavigateTo(constants.PAGE_VENDOR)
				})
			}
		}

		flex = flex.
			SetDirection(tview.FlexRow).
			AddItem(list, 0, 1, true).
			SetFullScreen(true)

		flex.SetBorder(true).SetBorderPadding(1, 1, 5, 5).SetTitle(" store ")

		return flex
	}, nil, onClose)
}

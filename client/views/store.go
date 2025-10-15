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
			for _, vendor := range resp.Vendors {
				list.AddItem("> "+vendor.GetVendorName(), vendor.GetVendorWareDescription()+"\r\n", 0, func() {
					c.AppStorage[constants.SELECTED_VENDOR_ID] = vendor.GetVendorId()
					c.AppStorage[constants.SELECTED_VENDOR_NAME] = vendor.GetVendorName()
					c.NavigateTo(constants.PAGE_VENDOR)
				})
			}
		}
		list = list.SetWrapAround(true)

		flex = flex.
			SetDirection(tview.FlexRow).
			AddItem(tview.NewTextView().SetText("available vendors: "), 2, 1, false).
			AddItem(list, 0, 1, true).
			AddItem(tview.NewTextView().SetText("weapons, armor, consumables, vanities..."), 1, 1, false).
			SetFullScreen(true)

		flex.SetBorder(true).SetBorderPadding(1, 1, 5, 5).SetTitle(" store ")

		return flex
	}, nil, onClose)
}

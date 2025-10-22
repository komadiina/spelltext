package views

import (
	"github.com/komadiina/spelltext/client/constants"
	types "github.com/komadiina/spelltext/client/types"
	"github.com/komadiina/spelltext/client/utils"
	pbRepo "github.com/komadiina/spelltext/proto/repo"
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
			m := utils.CreateModal("oops... ", "error: "+err.Error(), c, func() {
				c.NavigateTo(constants.PAGE_MAINMENU)
				c.App.SetRoot(c.PageManager.Pages, true).EnableMouse(true)
			})

			c.App.SetRoot(m, true).EnableMouse(true)
		}

		if len(resp.Vendors) == 0 {
			list.AddItem("woah...", "stock is empty. at the moment.", 'e', func() { onClose() })
		} else {
			for _, vendor := range resp.Vendors {
				list.AddItem("> "+vendor.GetVendorName(), vendor.GetVendorWareDescription()+"\r\n", 0, func() {
					c.Storage.SelectedVendor = &pbRepo.Vendor{
						Id:   vendor.VendorId,
						Name: vendor.VendorName,
					}
					c.NavigateTo(constants.PAGE_VENDOR)
				})
			}
		}
		list = list.SetWrapAround(true)

		guide := utils.CreateGuide([]*types.UnusableHotkey{
			{Key: "↑/↓", Desc: "navigate"},
			{Key: "enter", Desc: "select"},
			{Key: "esc", Desc: "back"},
		})

		flex = flex.
			SetDirection(tview.FlexRow).
			AddItem(tview.NewTextView().SetText("available vendors: "), 2, 1, false).
			AddItem(list, 0, 1, true).
			AddItem(tview.NewTextView().SetText("weapons, armor, consumables, vanities..."), 1, 1, false).
			AddItem(nil, 1, 1, false).
			AddItem(guide, 3, 1, false).
			SetFullScreen(true)

		flex.SetBorder(true).SetBorderPadding(1, 1, 5, 5).SetTitle(" [::b]store[::-] ")

		return flex
	}, nil, onClose)
}

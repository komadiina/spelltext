package views

import (
	types "github.com/komadiina/spelltext/client/types"
	pb "github.com/komadiina/spelltext/proto/store"
	"github.com/rivo/tview"
)

func AddStorePage(c *types.SpelltextClient) {
	onClose := func() {}

	c.PageManager.RegisterFactory(STORE_PAGE, func() tview.Primitive {
		flex := tview.NewFlex()
		list := tview.NewList()

		resp, err := c.Clients.StoreClient.ListItems(*c.Context, &pb.StoreListItemRequest{ItemType: 1})
		if err != nil {
			list.AddItem("no items fetched :(", "check storeserver logs", 'e', func() { onClose() })
		} else {
			for idx, item := range resp.Items {
				list.AddItem(item.GetName(), item.GetDescription(), rune(idx), func() {})
			}
		}

		flex.
			SetDirection(tview.FlexRow).
			AddItem(list, 0, 1, true)

		return flex
	}, nil, onClose)
}

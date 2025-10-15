package functions

import (
	"github.com/komadiina/spelltext/client/types"
	pbArmory "github.com/komadiina/spelltext/proto/armory"
	pbStore "github.com/komadiina/spelltext/proto/store"
)

func BuyItems(basket []*pbStore.Item, char *pbArmory.TCharacter, c *types.SpelltextClient) error {
	itemIds := make([]uint64, 0, len(basket))
	for _, item := range basket {
		itemIds = append(itemIds, item.GetId())
	}

	_, err := c.Clients.StoreClient.BuyItems(*c.Context, &pbStore.BuyItemRequest{CharacterId: char.GetId(), ItemIds: itemIds})
	return err
}

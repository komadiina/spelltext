package functions

import (
	"github.com/komadiina/spelltext/client/types"
	"github.com/komadiina/spelltext/client/utils"
	pbRepo "github.com/komadiina/spelltext/proto/repo"
	pbStore "github.com/komadiina/spelltext/proto/store"
)

func BuyItems(basket []*pbRepo.Item, char *pbRepo.Character, c *types.SpelltextClient) error {
	itemIds := make([]uint64, 0, len(basket))
	var cost int64 = 0
	for _, item := range basket {
		itemIds = append(itemIds, item.GetId())
		cost += int64(item.ItemTemplate.GoldPrice)
	}

	_, err := c.Clients.StoreClient.BuyItems(*c.Context, &pbStore.BuyItemRequest{CharacterId: char.GetCharacterId(), ItemIds: itemIds})
	if err == nil {
		utils.UpdateCharacterFunc(
			c.Storage.SelectedCharacter,
			c,
			func(t *pbRepo.Character) *pbRepo.Character {
				char := c.Storage.SelectedCharacter
				char.Gold = uint64(int64(char.Gold) - cost)
				return char
			},
		)
	}

	return err
}

package functions

import (
	"github.com/komadiina/spelltext/client/constants"
	"github.com/komadiina/spelltext/client/types"
	pbArmory "github.com/komadiina/spelltext/proto/armory"
	pbInventory "github.com/komadiina/spelltext/proto/inventory"
)

func GetBackpackItems(c *types.SpelltextClient) *pbInventory.ListBackpackItemsResponse {
	char := c.AppStorage[constants.SELECTED_CHARACTER].(*pbArmory.TCharacter)
	req := &pbInventory.ListBackpackItemsRequest{CharacterId: char.GetId()}
	resp, err := c.Clients.InventoryClient.ListBackpackItems(*c.Context, req)
	if err != nil {
		c.Logger.Error(err)
		return nil
	}

	return resp
}
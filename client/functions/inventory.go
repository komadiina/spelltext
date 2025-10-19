package functions

import (
	"fmt"
	"strings"

	"github.com/komadiina/spelltext/client/constants"
	"github.com/komadiina/spelltext/client/types"
	pbInventory "github.com/komadiina/spelltext/proto/inventory"
	pbRepo "github.com/komadiina/spelltext/proto/repo"
	"github.com/rivo/tview"
)

func GetBackpackItems(c *types.SpelltextClient) *pbInventory.ListBackpackItemsResponse {
	char := c.AppStorage[constants.SELECTED_CHARACTER].(*pbRepo.Character)
	req := &pbInventory.ListBackpackItemsRequest{CharacterId: char.GetCharacterId()}
	resp, err := c.Clients.InventoryClient.ListBackpackItems(*c.Context, req)
	if err != nil {
		c.Logger.Error(err)
		return nil
	}

	return resp
}

func GetRepoItemName(item *pbRepo.Item) string {
	return strings.Trim(fmt.Sprintf("%s%s%s", item.GetPrefix()+" ", item.GetItemTemplate().GetName(), " "+item.GetSuffix()), " ")
}

func MakeInventoryTableRow(row int, item *pbRepo.Item, c *types.SpelltextClient, t *tview.Table) *tview.Table {
	t = setCell(
		t, row, 0, GetRepoItemName(item), constants.COLOR_NAME, true,
		true,
	)

	t = setCell(
		t, row, 1, fmt.Sprint(item.GetItemTemplate().GetGoldPrice())+"g", constants.COLOR_PRICE, true,
		false,
	)

	t = setCell(
		t, row, 2, fmt.Sprint(item.GetHealth()), constants.COLOR_HEALTH, true,
		false,
	)

	t = setCell(
		t, row, 3, fmt.Sprint(item.GetPower()), constants.COLOR_POWER, true,
		false,
	)

	t = setCell(
		t, row, 4, fmt.Sprint(item.GetStrength()), constants.COLOR_STRENGTH, true,
		false,
	)

	t = setCell(
		t, row, 5, fmt.Sprint(item.GetSpellpower()), constants.COLOR_SPELLPOWER, true,
		false,
	)

	t = setCell(
		t, row, 6, fmt.Sprint(item.GetBonusDamage()), constants.COLOR_DAMAGE, true,
		false,
	)

	t = setCell(
		t, row, 7, fmt.Sprint(item.GetBonusArmor()), constants.COLOR_ARMOR, true,
		false,
	)

	return t
}

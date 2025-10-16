package functions

import (
	"fmt"

	"github.com/komadiina/spelltext/client/constants"
	"github.com/komadiina/spelltext/client/types"
	pbArmory "github.com/komadiina/spelltext/proto/armory"
	pbGamba "github.com/komadiina/spelltext/proto/gamba"
	pbRepo "github.com/komadiina/spelltext/proto/repo"
	"github.com/rivo/tview"
)

func GetGambaChests(c *types.SpelltextClient) (*pbGamba.GetChestsResponse, error) {
	resp, err := c.Clients.GambaClient.GetChests(*c.Context, &pbGamba.GetChestsRequest{})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func MakeChestTableHeader(table *tview.Table) *tview.Table {
	table = table.SetCell(0, 0, &tview.TableCell{
		Text:          "Name",
		Color:         constants.COLOR_NAME,
		Align:         tview.AlignLeft,
		NotSelectable: true,
		Expansion:     1,
	})

	table = table.SetCell(0, 1, &tview.TableCell{
		Text:          " Price",
		Color:         constants.COLOR_PRICE,
		Align:         tview.AlignLeft,
		NotSelectable: true,
	})

	return table
}

func MakeChestTableRow(row int, chest *pbRepo.GambaChest, table *tview.Table) *tview.Table {
	table = setCell(
		table, row, 0, chest.GetName(), constants.COLOR_NAME, true,
		true,
	)

	table = setCell(
		table, row, 1, fmt.Sprintf("%dg", chest.GetGoldPrice()), constants.COLOR_PRICE, true,
		false,
	)

	return table
}

func OpenChest(chest *pbRepo.GambaChest, c *types.SpelltextClient) (*pbGamba.OpenChestResponse, error) {
	char := c.AppStorage[constants.SELECTED_CHARACTER].(*pbArmory.TCharacter)
	if char.Gold < chest.GoldPrice {
		err := fmt.Errorf("unable to open chest, insufficient balance: need=>%dg, have=%dg", chest.GetGoldPrice(), char.GetGold())
		c.Logger.Error(err)
		return nil, err
	}

	req := &pbGamba.OpenChestRequest{CharacterId: char.GetId(), ChestId: chest.GetId()}
	resp, err := c.Clients.GambaClient.OpenChest(*c.Context, req)
	if err != nil {
		c.Logger.Error(err)
		return nil, err
	}

	return resp, nil
}

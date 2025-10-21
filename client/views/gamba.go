package views

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/komadiina/spelltext/client/constants"
	"github.com/komadiina/spelltext/client/functions"
	types "github.com/komadiina/spelltext/client/types"
	"github.com/komadiina/spelltext/client/utils"
	pbRepo "github.com/komadiina/spelltext/proto/repo"
	"github.com/rivo/tview"
)

func AddGambaPage(c *types.SpelltextClient) {
	onClose := func() {}

	c.PageManager.RegisterFactory(constants.PAGE_GAMBA, func() tview.Primitive {
		flex := tview.NewFlex().
			SetDirection(tview.FlexRow).
			SetFullScreen(true)
		flex.SetBorder(true).SetBorderPadding(1, 1, 5, 5).SetTitle(" [::b]gamba[::-] ")

		if c.AppStorage[constants.SELECTED_CHARACTER] == nil {
			f := tview.NewFlex().SetFullScreen(true)
			f.SetBorder(true).SetBorderPadding(1, 1, 5, 5).SetTitle(" hello? ")

			return f.AddItem(tview.NewTextView().
				SetText("no character selected. select a character from the character page, and come back... dummy"), 0, 1, false)
		}

		resp, err := functions.GetGambaChests(c)
		if err != nil {
			c.Logger.Error(err)
			flex.AddItem(tview.NewTextView().SetText("failed to fetch available chests. please try again later."), 0, 1, false)
		} else {
			chests := resp.GetChests()

			if len(chests) == 0 {
				flex.AddItem(tview.NewTextView().SetText("no chests exists."), 0, 1, false)
			} else {
				table := tview.NewTable().SetSeparator('|')
				table.SetBorder(true)
				table = functions.MakeChestTableHeader(table)

				availableGold := tview.NewTextView().
					SetText(
						fmt.Sprintf(
							"available gold: %dg",
							c.AppStorage[constants.SELECTED_CHARACTER].(*pbRepo.Character).GetGold(),
						))

				availableGold.SetBorder(true).SetBorderPadding(1, 1, 2, 2)

				reward := tview.NewTextView().
					SetDynamicColors(true).
					SetText("open a chest to get rewards..")

				reward.SetBorder(true).SetBorderPadding(1, 1, 2, 2)

				for idx, chest := range chests {
					table = functions.MakeChestTableRow(idx+1, chest, table)
				}

				table.
					Select(1, 0).
					SetFixed(1, 0).
					SetDoneFunc(func(key tcell.Key) {
						if key == tcell.KeyEnter {
							table.SetSelectable(true, false)
						}
					}).
					SetSelectedFunc(func(row, column int) {
						res, err := functions.OpenChest(resp.Chests[row-1], c)
						if err != nil {
							c.Logger.Error(err)
						} else {
							c.Logger.Info(res)
							availableGold = utils.UpdateGold(availableGold, "available gold: %dg", -int64(resp.Chests[row-1].GoldPrice), c)
							reward.SetText(
								fmt.Sprintf(
									`you won: [yellow]%s[""]! was it worth it?`,
									utils.GetFullItemName(res.GetItem()),
								))
						}

					})

				table.SetEvaluateAllRows(true)
				table.SetBorderPadding(1, 1, 5, 5)

				flex.
					AddItem(availableGold, 5, 1, false).
					AddItem(reward, 5, 1, false).
					AddItem(table, 0, 1, true)
			}
		}

		return flex
	}, nil, onClose)
}

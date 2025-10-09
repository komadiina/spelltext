package functions

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/komadiina/spelltext/client/constants"
	"github.com/komadiina/spelltext/client/utils"
	pb "github.com/komadiina/spelltext/proto/store"
	"github.com/rivo/tview"
)

func setCell(table *tview.Table, row int, column int, content string, color tcell.Color, selectable bool, alignLeft bool) *tview.Table {
	var alignment int
	if alignLeft {
		alignment = tview.AlignLeft
	} else {
		alignment = tview.AlignRight
	}

	table.SetCell(row, column,
		tview.NewTableCell(content).
			SetTextColor(color).
			SetAlign(alignment).
			SetSelectable(selectable),
	)

	return table
}

func MakeVendorTableRow(row int, item *pb.Item, table *tview.Table) *tview.Table {
	table = setCell(
		table, row, 0, utils.GetItemName(item), constants.COLOR_NAME, true,
		true,
	)

	table = setCell(
		table, row, 1,
		"0.00", // TODO: ADD PRICE FIELD
		constants.COLOR_PRICE, true, false,
	)

	table = setCell(
		table, row, 2, fmt.Sprint(item.GetHealthPoints()), constants.COLOR_HEALTH, true,
		false,
	)

	table = setCell(
		table, row, 3, fmt.Sprint(item.GetPowerPoints()), constants.COLOR_POWER, true,
		false,
	)

	table = setCell(
		table, row, 4, fmt.Sprint(item.GetStrengthPoints()), constants.COLOR_STRENGTH, true,
		false,
	)

	table = setCell(
		table, row, 5, fmt.Sprint(item.GetSpellpowerPoints()), constants.COLOR_SPELLPOWER, true,
		false,
	)

	table = setCell(
		table, row, 6, "0", constants.COLOR_DAMAGE, true, // TODO: DAMAGE IN SQL
		false,
	)

	return table
}

func MakeVendorTableHeader(table *tview.Table) *tview.Table {
	table = table.SetCell(0, 0, &tview.TableCell{
		Text:          "Name",
		Color:         constants.COLOR_NAME,
		Align:         tview.AlignLeft,
		NotSelectable: true,
		Expansion:     1,
	})

	table = table.SetCell(0, 1, &tview.TableCell{
		Text:          "Price",
		Color:         constants.COLOR_PRICE,
		Align:         tview.AlignLeft,
		NotSelectable: true,
	})

	table = table.SetCell(0, 2, &tview.TableCell{
		Text:          "Health",
		Color:         constants.COLOR_HEALTH,
		Align:         tview.AlignLeft,
		NotSelectable: true,
	})

	table = table.SetCell(0, 3, &tview.TableCell{
		Text:          "Power",
		Color:         constants.COLOR_POWER,
		Align:         tview.AlignLeft,
		NotSelectable: true,
	})

	table = table.SetCell(0, 4, &tview.TableCell{
		Text:          "Strength",
		Color:         constants.COLOR_STRENGTH,
		Align:         tview.AlignRight,
		NotSelectable: true,
	})

	table = table.SetCell(0, 5, &tview.TableCell{
		Text:          "Spellpower",
		Color:         constants.COLOR_SPELLPOWER,
		Align:         tview.AlignLeft,
		NotSelectable: true,
	})

	table = table.SetCell(0, 6, &tview.TableCell{
		Text:          "DPS",
		Color:         constants.COLOR_DAMAGE,
		Align:         tview.AlignLeft,
		NotSelectable: true,
	})

	return table
}

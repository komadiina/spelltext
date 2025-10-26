package views

import (
	"cmp"
	"fmt"
	"maps"
	"slices"

	"github.com/gdamore/tcell/v2"
	"github.com/komadiina/spelltext/client/constants"
	"github.com/komadiina/spelltext/client/functions"
	types "github.com/komadiina/spelltext/client/types"
	"github.com/komadiina/spelltext/client/utils"
	pb "github.com/komadiina/spelltext/proto/char"
	pbRepo "github.com/komadiina/spelltext/proto/repo"
	"github.com/rivo/tview"
)

type CharacterDetailsView struct {
	Name     *tview.TextView
	Level    *tview.TextView
	Class    *tview.TextView
	Currency *tview.TextView
}

type CharacterStatsView struct {
	HealthPoints *tview.TextView
	PowerPoints  *tview.TextView
	Strength     *tview.TextView
	Spellpower   *tview.TextView
	Armor        *tview.TextView
	Damage       *tview.TextView
}

var statsChanged bool = false
var onClose = func(*types.SpelltextClient) {}
var equipmentPane *tview.Flex = nil
var detailsPane *tview.Flex = nil

func (d *CharacterDetailsView) Update(c *pbRepo.Character) {
	d.Name.SetText(fmt.Sprintf("name: %s", c.CharacterName))
	d.Level.SetText(fmt.Sprintf(`level: [blue]%d[::-][white] (%d xp needed to level up)`, c.Level, constants.XP_CAP(c.Level)-c.Experience))
	d.Class.SetText(fmt.Sprintf(`class: %s`, c.Hero.Name))
	d.Currency.SetText(fmt.Sprintf(`[yellow]%dg[""] | [orange]%dt[""]`, c.Gold, c.Tokens))
}

func RenderCharactersList(
	details *CharacterDetailsView,
	chars *pb.ListCharactersResponse,
	stored []*pbRepo.Character,
	selected *tview.TextView,
	c *types.SpelltextClient,
) *tview.List {
	characters := tview.NewList()
	characters.SetTitleColor(tcell.ColorCadetBlue)
	for _, character := range chars.Characters {
		stored = append(stored, character)

		characters.AddItem("-> "+character.CharacterName, "", 0, func() {
			c.Storage.SelectedCharacter = character

			mod := tview.NewModal().
				SetText(fmt.Sprintf("character: %s", character.CharacterName)).
				AddButtons([]string{"select", "delete", "cancel"})

			mod.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				switch buttonIndex {
				case 0: // select
					functions.SetSelectedCharacter(character, c)
					c.Storage.SelectedCharacter = character
					c.Storage.CharacterStats = functions.CalculateStats(functions.GetEquippedItems(c), c) // recalculate stats
					c.App.SetRoot(c.PageManager.Pages, true).EnableMouse(true)
					selected.SetText(fmt.Sprintf(`+++ currently selected: [orange]%s[""] (lv. %d %s)`, character.CharacterName, character.Level, character.Hero.Name))
					return
				case 1: // delete
					functions.DeleteCharacter(character, c)
					c.App.SetRoot(c.PageManager.Pages, true).EnableMouse(true)
					c.PageManager.Pop()
					c.NavigateTo(constants.PAGE_CHARACTER)
					return
				case 2:
					c.App.SetRoot(c.PageManager.Pages, true).EnableMouse(true)
					return
				}
			})

			mod.SetBackgroundColor(tcell.ColorDarkSlateGrey)
			mod.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
				if event.Key() == tcell.KeyEscape {
					c.App.SetRoot(c.PageManager.Pages, true).EnableMouse(true)
					return nil
				}

				return event
			})
			c.App.SetRoot(mod, true).EnableMouse(true)
		})
	}

	details.Update(stored[0])

	characters.SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		details.Update(stored[index])
	})

	characters.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlA {
			c.PageManager.Push(constants.PAGE_CREATE_CHARACTER, false)
		}

		return event
	})

	characters.
		SetBorder(true).
		SetTitle(" [::b]characters[::-] ").
		SetTitleAlign(tview.AlignLeft).
		SetBorderPadding(1, 1, 2, 2)

	return characters
}

func RenderDetailsPane() (*tview.Flex, *CharacterDetailsView) {
	detailsPane = tview.NewFlex().SetDirection(tview.FlexRow)
	detailsPane.SetTitleColor(tcell.ColorCadetBlue)
	details := CharacterDetailsView{
		tview.NewTextView().SetDynamicColors(true),
		tview.NewTextView().SetDynamicColors(true),
		tview.NewTextView().SetDynamicColors(true),
		tview.NewTextView().SetDynamicColors(true),
	}
	detailsPane.SetBorder(true).
		SetBorderPadding(0, 1, 1, 1).
		SetTitle(" [::b]character details[::-] ").
		SetTitleAlign(tview.AlignLeft)

	detailsPane = detailsPane.
		AddItem(details.Name, 1, 1, false).
		AddItem(details.Level, 1, 1, false).
		AddItem(details.Class, 1, 1, false).
		AddItem(details.Currency, 1, 1, false)

	return detailsPane, &details
}

func RenderGuide() *tview.Flex {
	return utils.CreateGuide([]*types.UnusableHotkey{
		{Key: "ctrl+a", Desc: "new character"},
		{Key: "enter", Desc: "select"},
		{Key: "tab", Desc: "navigate"},
	}, true)
}

func RenderEquipmentPane(c *types.SpelltextClient, charSelected *pbRepo.Character) *tview.Flex {
	equipmentPane = tview.NewFlex().SetDirection(tview.FlexColumn)
	equipmentPane.SetTitleColor(tcell.ColorCadetBlue)

	equipmentPane.SetBorder(true).SetBorderPadding(1, 1, 2, 2).SetTitle(" [::b]equipment[::-] ")

	equipSlots := functions.GetEquipSlots(c)
	backpack := functions.GetBackpackItems(c).GetItemInstances() // :(
	grouped := functions.GroupItems(backpack, equipSlots)

	equipped := functions.GetEquippedItems(c)
	sortedKeys := slices.Collect(maps.Keys(grouped))
	slices.SortFunc(sortedKeys, func(a, b *pbRepo.EquipSlot) int {
		return cmp.Compare(a.Id, b.Id)
	})

	totals := RenderTotalsPane(c, equipped)
	bonusesPane := RenderBonusesPane(equipped)

	itemInfo := tview.NewTextView().SetDynamicColors(true)
	root := tview.NewTreeNode("quick inventory")
	tree := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)

	tree = RenderQuickInventoryTree(c, charSelected, equipSlots, tree, root, grouped, sortedKeys, itemInfo, totals)

	flexLeft := tview.NewFlex().SetDirection(tview.FlexRow)
	flexLeft.AddItem(itemInfo, 2, 1, false)
	flexLeft.AddItem(tree, 0, 2, true)
	equipmentPane.AddItem(flexLeft, 0, 2, true)

	bonusesPane.
		AddItem(nil, 0, 1, false).
		AddItem(totals, 4, 1, false)

	equipmentPane.AddItem(bonusesPane, 0, 3, false)
	return equipmentPane
}

func RenderQuickInventoryTree(
	c *types.SpelltextClient,
	charSelected *pbRepo.Character,
	equipSlots []*pbRepo.EquipSlot,
	tree *tview.TreeView,
	root *tview.TreeNode,
	grouped map[*pbRepo.EquipSlot][]*pbRepo.ItemInstance,
	sortedKeys []*pbRepo.EquipSlot,
	itemInfo *tview.TextView,
	totals *tview.Flex,
) *tview.TreeView {
	for _, slotKey := range sortedKeys {
		color := "#dedede"

		instances := grouped[slotKey]
		node := tview.NewTreeNode(fmt.Sprintf(`[%s]%s[""]`, color, slotKey.Name))
		for _, instance := range instances {
			color = "#afafaf"
			name := utils.GetFullItemName(instance.GetItem())

			child := tview.
				NewTreeNode(fmt.Sprintf(`[%s]%s[""]`, color, name)).
				SetSelectable(true)
			child.SetReference(instance)
			child.SetSelectedFunc(func() {
				functions.ToggleEquip(
					instance,
					c,
					charSelected,
					instance.Item.ItemTemplate.EquipSlot,
					true,
				)

				statsChanged = true

				// TODO: implement update equipped logic locally
				equipped := functions.GetEquippedItems(c)
				totals = RenderTotalsPane(c, equipped)
			})

			node.SetSelectable(true).SetSelectedFunc(func() {
				functions.ToggleEquip(instance, c, charSelected, instance.Item.ItemTemplate.EquipSlot, false)
			})

			node.AddChild(child)
		}

		root.AddChild(node)
	}

	tree.SetChangedFunc(func(node *tview.TreeNode) {
		ref := node.GetReference()
		if ref == nil {
			return
		}

		instance := ref.(*pbRepo.ItemInstance)
		itemInfo.SetText(utils.GetItemStats(instance.Item))
	})

	return tree
}

func RenderBonusesPane(equipped []*pbRepo.ItemInstance) *tview.Flex {
	bonusesPane := tview.NewFlex().SetDirection(tview.FlexRow)
	bonusesPane.SetTitleColor(tcell.ColorCadetBlue)
	bonusesPane.SetBorder(true).SetBorderPadding(1, 1, 2, 2).SetTitle(" [::b]bonuses[::-] ")

	slices.SortFunc(equipped, func(a, b *pbRepo.ItemInstance) int {
		return cmp.Compare(a.Item.ItemTemplate.EquipSlotId, b.Item.ItemTemplate.EquipSlotId)
	})

	for _, eqi := range equipped {
		es := "[" + eqi.GetItem().GetItemTemplate().GetEquipSlot().GetName() + "]" // escape tview color parsing
		tv := tview.NewTextView().SetDynamicColors(true)
		stats := utils.GetItemStats(eqi.Item)
		str := fmt.Sprintf(`[::b]%s[::-] %s%s[::i]%s[::-]`,
			tview.Escape(es), utils.GetFullItemName(eqi.GetItem()), "\n\t", stats,
		)

		tv.SetText(str).SetTextAlign(tview.AlignLeft)
		bonusesPane.AddItem(tv, 2, 1, false)
	}

	return bonusesPane
}

func RenderTotalsPane(c *types.SpelltextClient, equipped []*pbRepo.ItemInstance) *tview.Flex {
	totalsPane := tview.NewFlex().SetDirection(tview.FlexRow)
	totalsPane.SetBorder(true).SetBorderPadding(0, 0, 2, 2).SetTitle(" [::b]totals[::-] ")
	totalsPane.SetTitleColor(tcell.ColorCadetBlue)

	totals := &CharacterStatsView{
		HealthPoints: tview.NewTextView().SetDynamicColors(true),
		PowerPoints:  tview.NewTextView().SetDynamicColors(true),
		Strength:     tview.NewTextView().SetDynamicColors(true),
		Spellpower:   tview.NewTextView().SetDynamicColors(true),
		Armor:        tview.NewTextView().SetDynamicColors(true),
		Damage:       tview.NewTextView().SetDynamicColors(true),
	}

	// change CloserFunc here so i dont have to re-fetch nor re-store equipped items
	onClose = func(c *types.SpelltextClient) {
		cstats := functions.CalculateStats(functions.GetEquippedItems(c), c)
		c.Storage.CharacterStats = cstats
	}

	cstats := functions.CalculateStats(equipped, c)
	c.Storage.CharacterStats = cstats

	totals.HealthPoints.SetText(fmt.Sprintf(`[%s]HP[""]: %d`, constants.TEXT_COLOR_HEALTH, cstats.HealthPoints))
	totals.PowerPoints.SetText(fmt.Sprintf(`[%s]PWR[""]: %d`, constants.TEXT_COLOR_POWER, cstats.PowerPoints))
	totals.Strength.SetText(fmt.Sprintf(`[%s]STR[""]: %d`, constants.TEXT_COLOR_STRENGTH, cstats.Strength))
	totals.Spellpower.SetText(fmt.Sprintf(`[%s]SP[""]: %d`, constants.TEXT_COLOR_SPELLPOWER, cstats.Spellpower))
	totals.Armor.SetText(fmt.Sprintf(`[%s]ARM[""]: %d`, constants.TEXT_COLOR_ARMOR, cstats.Armor))
	totals.Damage.SetText(fmt.Sprintf(`[%s]DMG[""]: %d`, constants.TEXT_COLOR_DAMAGE, cstats.Damage))

	totalsPane.
		AddItem(
			tview.NewFlex().SetDirection(tview.FlexColumn).
				AddItem(totals.HealthPoints, 0, 1, false).
				AddItem(totals.PowerPoints, 0, 1, false).
				AddItem(totals.Strength, 0, 1, false), 0, 1, false).
		AddItem(
			tview.NewFlex().SetDirection(tview.FlexColumn).
				AddItem(totals.Spellpower, 0, 1, false).
				AddItem(totals.Armor, 0, 1, false).
				AddItem(totals.Damage, 0, 1, false), 0, 1, false)

	return totalsPane
}

func SetFlexInputHandler(flex *tview.Flex, equipmentPane *tview.Flex, characters *tview.List, c *types.SpelltextClient) {
	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTAB {
			focused := c.App.GetFocus()
			if focused == equipmentPane {
				c.App.SetFocus(characters)
			} else {
				c.App.SetFocus(equipmentPane)
			}
		} else if event.Key() == tcell.KeyCtrlA {
		}
		return event
	})
}

func AddCharacterPage(c *types.SpelltextClient) {
	c.PageManager.RegisterFactory(constants.PAGE_CHARACTER, func() tview.Primitive {
		flex := tview.NewFlex().SetDirection(tview.FlexRow)
		flex.SetBorder(true).SetBorderPadding(1, 1, 5, 5).SetTitle(" [::b]character[::-] ")
		flex = flex.SetFullScreen(true)

		header := tview.NewTextView().SetText("select one from available characters: ")

		var uid uint64 = 1 // TODO
		chars, err := functions.GetCharacters(uid, c)
		if err != nil {
			c.Logger.Error(err)
			return utils.GenerateErrorPage(c, err.Error())
		}

		charSelected := c.Storage.SelectedCharacter
		if charSelected == nil {
			charSelected = &pbRepo.Character{Hero: &pbRepo.Hero{Name: "none selected.."}}
		}
		selected := tview.NewTextView().
			SetText(
				fmt.Sprintf(`+++ currently selected: [orange]%s[""]!`, charSelected.CharacterName)).
			SetDynamicColors(true)

		detailsPane, details := RenderDetailsPane()

		characters := RenderCharactersList(details, chars, []*pbRepo.Character{}, selected, c)

		guide := RenderGuide()
		equipmentPane = RenderEquipmentPane(c, charSelected)
		SetFlexInputHandler(flex, equipmentPane, characters, c)

		flex.
			AddItem(detailsPane, 6, 1, false).
			AddItem(nil, 1, 1, false).
			AddItem(selected, 1, 1, false).
			AddItem(header, 1, 1, false).
			AddItem(nil, 1, 1, false).
			AddItem(characters, 0, 1, true).
			AddItem(nil, 1, 1, false).
			AddItem(equipmentPane, 0, 3, false).
			AddItem(guide, 3, 1, false)

		return flex
	}, nil, func() { onClose(c) })
}

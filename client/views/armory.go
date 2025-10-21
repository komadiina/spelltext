package views

import (
	"cmp"
	"fmt"
	"maps"
	"slices"
	"strings"

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

func (d *CharacterDetailsView) Update(c *pbRepo.Character) {
	d.Name.SetText(fmt.Sprintf("name: %s", c.CharacterName))
	d.Level.SetText(fmt.Sprintf(`level: [blue]%d[""]`, c.Level))
	d.Class.SetText(fmt.Sprintf(`class: %s`, c.Hero.Name))
	d.Currency.SetText(fmt.Sprintf(`[yellow]%dg[""] | [orange]%dt[""]`, c.Gold, c.Tokens))
}

func sumStats(inst *pbRepo.ItemInstance, cstats *types.CharacterStats) *types.CharacterStats {
	return &types.CharacterStats{
		HealthPoints: cstats.HealthPoints + inst.Item.Health,
		PowerPoints:  cstats.PowerPoints + inst.Item.Power,
		Strength:     cstats.Strength + inst.Item.Strength,
		Spellpower:   cstats.Spellpower + inst.Item.Spellpower,
		Armor:        cstats.Armor + inst.Item.BonusArmor,
		Damage:       cstats.Damage + inst.Item.BonusDamage,
	}
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
	detailsPane := tview.NewFlex().SetDirection(tview.FlexRow)
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
	guide := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(tview.NewTextView().SetText(" keymap legend: "), 0, 1, false)
	guide.SetBorder(true)

	back, length := utils.AddNavGuide("esc", "back")
	guide.AddItem(back, length, 1, false)

	add, length := utils.AddNavGuide("ctrl+a", "create new character")
	guide.AddItem(add, length, 1, false)

	enter, length := utils.AddNavGuide("enter", "character menu")
	guide.AddItem(enter, length, 1, false)
	return guide
}

func RenderEquipmentPane(c *types.SpelltextClient, charSelected *pbRepo.Character) *tview.Flex {
	equipmentPane := tview.NewFlex().SetDirection(tview.FlexColumn)
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

	root := tview.NewTreeNode("quick inventory")
	tree := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)
	tree = RenderQuickInventoryTree(c, charSelected, equipSlots, tree, root, grouped, sortedKeys)

	equipmentPane.AddItem(tree, 0, 2, true)

	totals := RenderTotalsPane(equipped)
	bonusesPane := RenderBonusesPane(equipped)

	bonusesPane.
		AddItem(nil, 0, 1, false).
		AddItem(totals, 3, 1, false)

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
			child.SetSelectedFunc(func() {
				functions.ToggleEquip(
					instance,
					c,
					charSelected,
					instance.Item.ItemTemplate.EquipSlot,
					true,
				)
				// TODO: refresh equipped items
			})
			node.AddChild(child)
		}

		root.AddChild(node)
	}

	return tree
}

func RenderBonusesPane(equipped []*pbRepo.ItemInstance) *tview.Flex {
	bonusesPane := tview.NewFlex().SetDirection(tview.FlexRow)
	bonusesPane.SetTitleColor(tcell.ColorCadetBlue)
	bonusesPane.SetBorder(true).SetBorderPadding(1, 1, 2, 2).SetTitle(" [::b]bonuses[::-] ")

	slices.SortFunc(equipped, func(a, b *pbRepo.ItemInstance) int {
		return cmp.Compare(a.Item.ItemTemplate.EquipSlotId, b.Item.ItemTemplate.EquipSlotId)
	})

	discardEmpty := func(eq *pbRepo.ItemInstance) string {
		sb := strings.Builder{}

		if eq.Item.Health != 0 {
			sgn := "+"
			if eq.Item.Health < 0 {
				sgn = ""
			}

			sb.WriteString(fmt.Sprintf(`[%s]%s%d HP[""], `, constants.TEXT_COLOR_HEALTH, sgn, eq.Item.Health))
		}

		if eq.Item.Power != 0 {
			sgn := "+"
			if eq.Item.Power < 0 {
				sgn = ""
			}

			sb.WriteString(fmt.Sprintf(`[%s]%s%d PWR[""], `, constants.TEXT_COLOR_POWER, sgn, eq.Item.Power))
		}

		if eq.Item.Strength != 0 {
			sgn := "+"
			if eq.Item.Strength < 0 {
				sgn = ""
			}

			sb.WriteString(fmt.Sprintf(`[%s]%s%d STR[""], `, constants.TEXT_COLOR_STRENGTH, sgn, eq.Item.Strength))
		}

		if eq.Item.Spellpower != 0 {
			sgn := "+"
			if eq.Item.Spellpower < 0 {
				sgn = ""
			}

			sb.WriteString(fmt.Sprintf(`[%s]%s%d SP[""], `, constants.TEXT_COLOR_SPELLPOWER, sgn, eq.Item.Spellpower))
		}

		if eq.Item.BonusDamage != 0 {
			sgn := "+"
			if eq.Item.BonusDamage < 0 {
				sgn = ""
			}

			sb.WriteString(fmt.Sprintf(`[%s]%s%d DMG[""], `, constants.TEXT_COLOR_DAMAGE, sgn, eq.Item.BonusDamage))
		}

		if eq.Item.BonusArmor != 0 {
			sgn := "+"
			if eq.Item.BonusArmor < 0 {
				sgn = ""
			}

			sb.WriteString(fmt.Sprintf(`[%s]%s%d ARM[""], `, constants.TEXT_COLOR_ARMOR, sgn, eq.Item.BonusArmor))
		}

		return sb.String()[:len(sb.String())-2]
	}

	for _, eqi := range equipped {
		es := "[" + eqi.GetItem().GetItemTemplate().GetEquipSlot().GetName() + "]" // escape tview color parsing
		tv := tview.NewTextView().SetDynamicColors(true)
		stats := discardEmpty(eqi)
		str := fmt.Sprintf(`[::b]%s[::-] %s%s[::i]%s[::-]`,
			tview.Escape(es), utils.GetFullItemName(eqi.GetItem()), "\n\t", stats,
		)

		tv.SetText(str).SetTextAlign(tview.AlignLeft)
		bonusesPane.AddItem(tv, 2, 1, false)
	}

	return bonusesPane
}

func RenderTotalsPane(equipped []*pbRepo.ItemInstance) *tview.Flex {
	totalsPane := tview.NewFlex().SetDirection(tview.FlexColumn)
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

	cstats := &types.CharacterStats{}
	for _, eqi := range equipped {
		cstats = sumStats(eqi, cstats)
	}

	totals.HealthPoints.SetText(fmt.Sprintf(`[%s]HP[""]: %d`, constants.TEXT_COLOR_HEALTH, cstats.HealthPoints))
	totals.PowerPoints.SetText(fmt.Sprintf(`[%s]PWR[""]: %d`, constants.TEXT_COLOR_POWER, cstats.PowerPoints))
	totals.Strength.SetText(fmt.Sprintf(`[%s]STR[""]: %d`, constants.TEXT_COLOR_STRENGTH, cstats.Strength))
	totals.Spellpower.SetText(fmt.Sprintf(`[%s]SP[""]: %d`, constants.TEXT_COLOR_SPELLPOWER, cstats.Spellpower))
	totals.Armor.SetText(fmt.Sprintf(`[%s]ARM[""]: %d`, constants.TEXT_COLOR_ARMOR, cstats.Armor))
	totals.Damage.SetText(fmt.Sprintf(`[%s]DMG[""]: %d`, constants.TEXT_COLOR_DAMAGE, cstats.Damage))

	totalsPane.
		AddItem(totals.HealthPoints, 0, 1, false).
		AddItem(totals.PowerPoints, 0, 1, false).
		AddItem(totals.Strength, 0, 1, false).
		AddItem(totals.Spellpower, 0, 1, false).
		AddItem(totals.Armor, 0, 1, false).
		AddItem(totals.Damage, 0, 1, false)

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
	onClose := func() {}

	c.PageManager.RegisterFactory(constants.PAGE_CHARACTER, func() tview.Primitive {
		flex := tview.NewFlex().SetDirection(tview.FlexRow)
		flex.SetBorder(true).SetBorderPadding(1, 1, 5, 5).SetTitle(" [::b]character[::-] ")
		flex = flex.SetFullScreen(true)

		header := tview.NewTextView().SetText("select one from available characters: ")

		var uid uint64 = 1 // TODO
		chars, err := functions.GetCharacters(uid, c)
		if err != nil {
			c.Logger.Error(err)
			chars = &pb.ListCharactersResponse{
				Characters: make([]*pbRepo.Character, 0),
			}
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
		equipmentPane := RenderEquipmentPane(c, charSelected)
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
	}, nil, onClose)
}

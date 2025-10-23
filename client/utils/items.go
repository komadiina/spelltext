package utils

import (
	"fmt"
	"strings"

	"github.com/komadiina/spelltext/client/constants"
	pbRepo "github.com/komadiina/spelltext/proto/repo"
)

func GetFullItemName(item *pbRepo.Item) string {
	return strings.Trim(fmt.Sprintf("%s %s %s", item.GetPrefix(), item.GetItemTemplate().GetName(), item.GetSuffix()), " ")
}

func GetItemName(item *pbRepo.Item) string {
	var prefix string = ""
	var suffix string = ""

	if len(item.GetPrefix()) == 0 {
		prefix = ""
	} else {
		prefix = item.GetPrefix() + " "
	}

	if len(item.GetSuffix()) == 0 {
		suffix = ""
	} else {
		suffix = " " + item.GetSuffix()
	}

	return fmt.Sprintf("%s%s%s", prefix, item.ItemTemplate.Name, suffix)
}

func GetItemStats(item *pbRepo.Item) string {
	sb := strings.Builder{}

	if item.Health != 0 {
		sgn := "+"
		if item.Health < 0 {
			sgn = ""
		}

		sb.WriteString(fmt.Sprintf(`[%s]%s%d HP[""], `, constants.TEXT_COLOR_HEALTH, sgn, item.Health))
	}

	if item.Power != 0 {
		sgn := "+"
		if item.Power < 0 {
			sgn = ""
		}

		sb.WriteString(fmt.Sprintf(`[%s]%s%d PWR[""], `, constants.TEXT_COLOR_POWER, sgn, item.Power))
	}

	if item.Strength != 0 {
		sgn := "+"
		if item.Strength < 0 {
			sgn = ""
		}

		sb.WriteString(fmt.Sprintf(`[%s]%s%d STR[""], `, constants.TEXT_COLOR_STRENGTH, sgn, item.Strength))
	}

	if item.Spellpower != 0 {
		sgn := "+"
		if item.Spellpower < 0 {
			sgn = ""
		}

		sb.WriteString(fmt.Sprintf(`[%s]%s%d SP[""], `, constants.TEXT_COLOR_SPELLPOWER, sgn, item.Spellpower))
	}

	if item.BonusDamage != 0 {
		sgn := "+"
		if item.BonusDamage < 0 {
			sgn = ""
		}

		sb.WriteString(fmt.Sprintf(`[%s]%s%d DMG[""], `, constants.TEXT_COLOR_DAMAGE, sgn, item.BonusDamage))
	}

	if item.BonusArmor != 0 {
		sgn := "+"
		if item.BonusArmor < 0 {
			sgn = ""
		}

		sb.WriteString(fmt.Sprintf(`[%s]%s%d ARM[""], `, constants.TEXT_COLOR_ARMOR, sgn, item.BonusArmor))
	}

	return sb.String()[:len(sb.String())-2]
}

package utils

import (
	"fmt"
	"strings"

	"github.com/komadiina/spelltext/client/constants"
	"github.com/komadiina/spelltext/client/types"
	pbRepo "github.com/komadiina/spelltext/proto/repo"
)

func GetFullNpcName(npc *pbRepo.Npc) string {
	return strings.Trim(fmt.Sprintf("%s %s %s", npc.GetPrefix(), npc.GetNpcTemplate().GetName(), npc.GetSuffix()), " ")
}

func PrintNpcDetails(npc *pbRepo.Npc) string {
	return fmt.Sprintf(`name: [white]%s[""]
level: [white]%d[""]
[white]base health:[""] [%s]%d[""]
[white]base damage:[""] [%s]%d[""]
[white]xp reward:[""] [%s]%d[""]
[white]gold reward:[""] [%s]%d[""]`,
		GetFullNpcName(npc),
		npc.Level,
		constants.TEXT_COLOR_HEALTH,
		int(float32(npc.NpcTemplate.HealthPoints)*npc.HealthMultiplier),
		constants.TEXT_COLOR_DAMAGE,
		int(float32(npc.NpcTemplate.BaseDamage)*npc.DamageMultiplier),
		constants.TEXT_COLOR_XP,
		npc.NpcTemplate.BaseXpReward,
		constants.TEXT_COLOR_GOLD,
		npc.NpcTemplate.GoldReward,
	)

}

func GetDisplayStatsPlayer(c *types.SpelltextClient) *[]uint64 {
	arr := make([]uint64, 0)
	arr = append(arr, uint64(c.Storage.CharacterStats.HealthPoints))
	arr = append(arr, uint64(c.Storage.CharacterStats.PowerPoints))
	return &arr
}

func GetDisplayStatsNpc(npc *pbRepo.Npc) *[]uint64 {
	arr := make([]uint64, 0)
	arr = append(arr,
		uint64(float32(npc.NpcTemplate.HealthPoints)*npc.HealthMultiplier),
		uint64(float32(npc.NpcTemplate.BaseDamage)*npc.DamageMultiplier),
	)
	return &arr
}

package utils

import (
	"fmt"
	"strings"

	"github.com/komadiina/spelltext/client/constants"
	pbRepo "github.com/komadiina/spelltext/proto/repo"
)

func GetFullNpcName(npc *pbRepo.Npc) string {
	return strings.Trim(fmt.Sprintf("%s %s %s", npc.GetPrefix(), npc.GetNpcTemplate().GetName(), npc.GetSuffix()), " ")
}

func PrintNpcDetails(npc *pbRepo.Npc) string {
	return fmt.Sprintf(`name: [white]%s[""]
level: [white]%d[""]
[white]base health:[""] [%s]%d[""]
[white]base damage:[""] [%s]%d[""]`,
		GetFullNpcName(npc),
		npc.Level,
		constants.TEXT_COLOR_HEALTH,
		int(float32(npc.NpcTemplate.HealthPoints)*npc.HealthMultiplier),
		constants.TEXT_COLOR_DAMAGE,
		int(float32(npc.NpcTemplate.BaseDamage)*npc.DamageMultiplier))

}

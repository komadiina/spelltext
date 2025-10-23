package functions

import (
	"strings"

	"github.com/komadiina/spelltext/client/constants"
	"github.com/komadiina/spelltext/client/types"
)

func RedrawHealthBar(npcState *types.NpcFightState) string {
	numBlocks := int64(float64(npcState.CurrentHealth) / float64(npcState.Npc.NpcTemplate.HealthPoints*uint64(npcState.Npc.HealthMultiplier)))

	sb := strings.Builder{}
	for i := 0; i < int(numBlocks); i++ {
		sb.WriteString(constants.CHARACTER_HEALTH)
	}

	return sb.String()
}

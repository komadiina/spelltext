package utils

import (
	"fmt"
	"strings"

	pbRepo "github.com/komadiina/spelltext/proto/repo"
)

func GetFullNpcName(npc *pbRepo.Npc) string {
	return strings.Trim(fmt.Sprintf("%s %s %s", npc.GetPrefix(), npc.GetNpcTemplate().GetName(), npc.GetSuffix()), " ")
}

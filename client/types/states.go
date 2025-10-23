package types

import (
	pbRepo "github.com/komadiina/spelltext/proto/repo"
)

type NpcFightState struct {
	Npc           *pbRepo.Npc
	CurrentHealth int64
}

package types

import (
	pbRepo "github.com/komadiina/spelltext/proto/repo"
)

type NpcFightState struct {
	Npc           *pbRepo.Npc
	CurrentHealth int64
}

type PlayerAbilities struct {
	available []*pbRepo.Ability
	unlocked  []*pbRepo.Ability
	upgraded  []*pbRepo.Ability
}

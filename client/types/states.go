package types

import (
	pbRepo "github.com/komadiina/spelltext/proto/repo"
)

type CbFightState struct {
	Npc *pbRepo.Npc

	NpcCurrentHealth    int64
	NpcCurrentPower     int64
	PlayerCurrentHealth int64
	PlayerCurrentPower  int64
}

type PlayerAbilities struct {
	available []*pbRepo.Ability
	unlocked  []*pbRepo.Ability
	upgraded  []*pbRepo.Ability
}

package functions

import (
	"math"

	pbRepo "github.com/komadiina/spelltext/proto/repo"
	"github.com/komadiina/spelltext/server/character/triggers"
)

func GetLevelXpCap(level uint64) uint64 {
	// equiv to: f(l) = 125l + 66l^2 ~= 66l(2l + 1)
	return uint64(float64(level)*100*1.25 + math.Pow(float64(level), 2.0)*66.0)
}

func AddXp(c *pbRepo.Character, xp uint64) *pbRepo.Character {
	xpCap := GetLevelXpCap(c.Level)
	if c.Experience+xp >= xpCap {
		c.Experience = c.Experience + xp - xpCap
		triggers.OnLevelUp(c)
	} else {
		c.Experience += xp
	}

	return c
}

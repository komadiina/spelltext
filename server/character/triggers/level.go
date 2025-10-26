package triggers

import (
	pbRepo "github.com/komadiina/spelltext/proto/repo"
)

func OnLevelUp(c *pbRepo.Character) {
	c.PointsHealth += uint64(c.Hero.HealthPerLevel)
	c.PointsPower += uint64(c.Hero.PowerPerLevel)
	c.PointsStrength += uint64(c.Hero.StrengthPerLevel)
	c.PointsSpellpower += uint64(c.Hero.SpellpowerPerLevel)
	c.UnspentPoints += 1
	c.Level += 1
}

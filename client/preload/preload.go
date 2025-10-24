package preload

import (
	"sync"

	"github.com/komadiina/spelltext/client/functions"
	"github.com/komadiina/spelltext/client/types"
)

var once sync.Once

func Preload(c *types.SpelltextClient) {
	once.Do(func() {
		InitializeClientStats(c)
	})
}

func InitializeClientStats(c *types.SpelltextClient) {
	cstats := functions.CalculateStats(functions.GetEquippedItems(c), c)
	c.Storage.CharacterStats = cstats

	c.Logger.Debug("[startup] initialized character stats.")
	c.Logger.Debug(c.Storage.CharacterStats)
}

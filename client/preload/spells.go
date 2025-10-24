package preload

import (
	"sync"

	"github.com/komadiina/spelltext/client/types"
)

var spOnce sync.Once

func InitSpellProcs(c *types.SpelltextClient) {
	spOnce.Do(func() {
		
	})
}


func CreateSpellProc(c *types.SpelltextClient, spellId uint64, proc *types.SpellProc) {
	c.Storage.SpellProcs[spellId] = proc
}
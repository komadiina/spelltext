package functions

import (
	"github.com/komadiina/spelltext/client/types"
	pbCombat "github.com/komadiina/spelltext/proto/combat"
	pbRepo "github.com/komadiina/spelltext/proto/repo"
)

func GetAvailableNpcs(c *types.SpelltextClient) []*pbRepo.Npc {
	resp, err := c.Clients.CombatClient.ListNpcs(*c.Context, &pbCombat.ListNpcsRequest{
		CharacterId: c.Storage.CurrentUser.Id,
	})

	if err != nil {
		c.Logger.Error(err)
		c.PageManager.Pop()
		return nil
	}

	return resp.GetNpcs()
}

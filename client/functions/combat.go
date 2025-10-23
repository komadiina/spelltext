package functions

import (
	"github.com/komadiina/spelltext/client/types"
	pbBuild "github.com/komadiina/spelltext/proto/build"
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

func GetPlayerSpells(c *types.SpelltextClient) []*pbRepo.PlayerAbilityTree {
	resp, err := c.Clients.BuildClient.ListAbilities(*c.Context, &pbBuild.ListAbilitiesRequest{Character: c.Storage.SelectedCharacter})
	if err != nil {
		return nil
	}

	return resp.Upgraded // alias, 'unlocked'
}

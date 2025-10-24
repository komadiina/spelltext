package functions

import (
	// "github.com/komadiina/spelltext/client/constants"
	"github.com/komadiina/spelltext/client/types"
	pbChar "github.com/komadiina/spelltext/proto/char"
	pbRepo "github.com/komadiina/spelltext/proto/repo"
)

func CreateCharacter(char *pbRepo.Character, c *types.SpelltextClient) error {
	req := &pbChar.CreateCharacterRequest{
		Hero:   char.GetHero(),
		Name:   char.GetCharacterName(),
		UserId: c.Storage.CurrentUser.GetId(),
	}

	c.Logger.Debug(req)

	resp, err := c.Clients.CharacterClient.CreateCharacter(*c.Context, req)
	if err != nil {
		c.Logger.Error(err)
		return err
	} else {
		err := SetSelectedCharacter(resp.GetCharacter(), c)
		if err != nil {
			c.Logger.Error(err)
			return err
		}

		c.PageManager.Pop()
	}

	return nil
}

func ListHeroes(c *types.SpelltextClient) []*pbRepo.Hero {
	resp, err := c.Clients.CharacterClient.ListHeroes(*c.Context, &pbChar.ListHeroesRequest{})
	if err != nil {
		c.Logger.Error(err)
		return nil
	}

	return resp.GetHeroes()
}

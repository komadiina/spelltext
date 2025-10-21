package functions

import (
	"github.com/komadiina/spelltext/client/constants"
	"github.com/komadiina/spelltext/client/types"
	pbAuth "github.com/komadiina/spelltext/proto/auth"
	pbChar "github.com/komadiina/spelltext/proto/char"
)

func LoginUser(c *types.SpelltextClient, username, password string) {
	req := &pbAuth.LoginRequest{Username: username, PasswordHash: "TODO"}
	resp, err := c.Clients.AuthClient.Login(*c.Context, req)
	if err != nil {
		c.Logger.Error(err)
	}

	c.AppStorage[constants.CURRENT_USER] = resp.GetUser()
	c.AppStorage[constants.SELECTED_CHARACTER] = resp.GetCharacter()
}

func SetLastSelectedCharacter(c *types.SpelltextClient) {
	req := &pbChar.GetLastSelectedCharacterRequest{
		Username: c.AppStorage[constants.CURRENT_USER_NAME].(string),
	}

	resp, err := c.Clients.CharacterClient.GetLastSelectedCharacter(*c.Context, req)
	if err != nil {
		c.Logger.Error(err)
	}

	c.AppStorage[constants.SELECTED_CHARACTER] = resp.GetCharacter()
}

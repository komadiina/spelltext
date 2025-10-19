package functions

import (
	"fmt"

	"github.com/komadiina/spelltext/client/constants"
	"github.com/komadiina/spelltext/client/types"
	pb "github.com/komadiina/spelltext/proto/armory"
	pbRepo "github.com/komadiina/spelltext/proto/repo"
)

func GetCharacters(uid uint64, c *types.SpelltextClient) (*pb.ListCharactersResponse, error) {
	resp, err := c.Clients.CharacterClient.ListCharacters(*c.Context, &pb.ListCharactersRequest{Username: c.User.Username})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func SetSelectedCharacter(char *pbRepo.Character, c *types.SpelltextClient) error {
	if char == nil {
		return fmt.Errorf("cant set c.AppStorage[%v], character is nil.", constants.SELECTED_CHARACTER)
	}

	req := &pb.SetSelectedCharacterRequest{
		CharacterId: char.GetCharacterId(),
		UserId:      c.AppStorage[constants.CURRENT_USER].(*pbRepo.User).GetId(),
	}

	_, err := c.Clients.CharacterClient.SetSelectedCharacter(*c.Context, req)

	if err != nil {
		c.Logger.Error(err)
		return err
	}

	c.AppStorage[constants.SELECTED_CHARACTER] = char
	return nil
}

func DeleteCharacter(char *pbRepo.Character, c *types.SpelltextClient) error {
	_, err := c.Clients.CharacterClient.DeleteCharacter(*c.Context, &pb.DeleteCharacterRequest{CharacterId: char.GetCharacterId()})

	if err != nil {
		c.Logger.Error(err)
		return fmt.Errorf("error=%v", err)
	}

	return nil
}

func RefreshCharacter(char *pbRepo.Character, c *types.SpelltextClient) error {
	resp, err := GetCharacters(c.AppStorage[constants.CURRENT_USER_ID].(uint64), c)

	for _, character := range resp.Characters {
		if character.GetCharacterId() == char.GetCharacterId() {
			return SetSelectedCharacter(character, c)
		}
	}

	return err
}

func GetEquippedItems(c *types.SpelltextClient) []*pbRepo.Item {
	req := &pb.GetEquippedItemsRequest{
		CharacterId: c.AppStorage[constants.SELECTED_CHARACTER].(*pbRepo.Character).GetCharacterId(),
	}

	resp, err := c.Clients.CharacterClient.GetEquippedItems(*c.Context, req)
	if err != nil {
		c.Logger.Error(err)
		return nil
	}

	return resp.GetItems()
}

package functions

import (
	"fmt"

	"github.com/komadiina/spelltext/client/constants"
	"github.com/komadiina/spelltext/client/types"
	pb "github.com/komadiina/spelltext/proto/armory"
)

func GetCharacters(uid uint64, c *types.SpelltextClient) (*pb.ListCharactersResponse, error) {
	resp, err := c.Clients.CharacterClient.ListCharacters(*c.Context, &pb.ListCharactersRequest{Username: c.User.Username})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func SetSelectedCharacter(char *pb.TCharacter, c *types.SpelltextClient) error {
	if char == nil {
		return fmt.Errorf("cant set c.AppStorage[%v], character is nil.", constants.SELECTED_CHARACTER)
	}

	c.AppStorage[constants.SELECTED_CHARACTER] = char
	return nil
}

func DeleteCharacter(char *pb.TCharacter, c *types.SpelltextClient) error {
	resp, err := c.Clients.CharacterClient.DeleteCharacter(*c.Context, &pb.DeleteCharacterRequest{CharacterId: char.GetId()})
	if err != nil {
		return fmt.Errorf("%s: error=%v", resp.GetMessage(), err)
	} else {
		return nil
	}
}

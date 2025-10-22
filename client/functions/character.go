package functions

import (
	"fmt"

	"github.com/komadiina/spelltext/client/types"
	pb "github.com/komadiina/spelltext/proto/char"
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
		return fmt.Errorf("cant set c.Storage.SelectedCharacter, character is nil.")
	}

	req := &pb.SetSelectedCharacterRequest{
		CharacterId: char.GetCharacterId(),
		UserId:      c.Storage.CurrentUser.GetId(),
	}

	_, err := c.Clients.CharacterClient.SetSelectedCharacter(*c.Context, req)

	if err != nil {
		c.Logger.Error(err)
		return err
	}

	c.Storage.SelectedCharacter = char
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
	resp, err := GetCharacters(c.Storage.CurrentUser.Id, c)

	for _, character := range resp.Characters {
		if character.GetCharacterId() == char.GetCharacterId() {
			return SetSelectedCharacter(character, c)
		}
	}

	return err
}

func GetEquippedItems(c *types.SpelltextClient) []*pbRepo.ItemInstance {
	req := &pb.GetEquippedItemsRequest{
		CharacterId: c.Storage.SelectedCharacter.GetCharacterId(),
	}

	resp, err := c.Clients.CharacterClient.GetEquippedItems(*c.Context, req)
	if err != nil {
		c.Logger.Error(err)
		return nil
	}

	return resp.GetItemInstances()
}

func GetEquipSlots(c *types.SpelltextClient) []*pbRepo.EquipSlot {
	if c.Storage.EquipSlots != nil || len(c.Storage.EquipSlots) > 0 {
		return c.Storage.EquipSlots
	}

	resp, err := c.Clients.CharacterClient.GetEquipSlots(*c.Context, &pb.GetEquipSlotsRequest{})
	if err != nil {
		c.Logger.Error(err)
		return nil
	}

	c.Storage.EquipSlots = resp.GetSlots()
	return resp.GetSlots()
}

func GroupItems(instances []*pbRepo.ItemInstance, slots []*pbRepo.EquipSlot) map[*pbRepo.EquipSlot][]*pbRepo.ItemInstance {
	m := make(map[*pbRepo.EquipSlot][]*pbRepo.ItemInstance)
	// slices.SortFunc(slots, func(a, b *pbRepo.EquipSlot) int { return cmp.Compare(a.GetId(), b.GetId()) })
	for _, instance := range instances {
		for _, slot := range slots {
			if instance.GetItem().GetItemTemplate().GetEquipSlot().GetId() == slot.GetId() {
				m[slot] = append(m[slot], instance)
			}
		}
	}

	return m
}

func ToggleEquip(item *pbRepo.ItemInstance, c *types.SpelltextClient, char *pbRepo.Character, equipSlot *pbRepo.EquipSlot, shouldEquip bool) error {
	req := &pb.ToggleEquipRequest{
		CharacterId:    char.GetCharacterId(),
		ItemInstanceId: item.ItemInstanceId,
		EquipSlotId:    equipSlot.GetId(),
		ShouldEquip:    shouldEquip,
	}

	_, err := c.Clients.CharacterClient.ToggleEquip(*c.Context, req)

	if err != nil {
		c.Logger.Error(err)
		return err
	}

	return nil
}

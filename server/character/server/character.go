package server

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	sq "github.com/Masterminds/squirrel"

	pb "github.com/komadiina/spelltext/proto/char"
	pbHealth "github.com/komadiina/spelltext/proto/health"
	pbRepo "github.com/komadiina/spelltext/proto/repo"
	"github.com/komadiina/spelltext/server/character/config"
	"github.com/komadiina/spelltext/server/character/functions"
	monitor "github.com/komadiina/spelltext/utils"
	"github.com/komadiina/spelltext/utils/singleton/logging"
)

type CharacterService struct {
	pb.UnimplementedCharacterServer
	monitor.HealthCheckable
	DbPool *pgxpool.Pool
	Config *config.Config
	Logger *logging.Logger
}

func (s *CharacterService) Check(ctx context.Context, req *pbHealth.HealthCheckRequest) (*pbHealth.HealthCheckResponse, error) {
	return &pbHealth.HealthCheckResponse{Status: pbHealth.HealthCheckResponse_SERVING}, nil
}

func (s *CharacterService) ListHeroes(ctx context.Context, req *pb.ListHeroesRequest) (*pb.ListHeroesResponse, error) {
	sql, _, err := sq.Select("*").From("heroes").ToSql()

	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	rows, err := s.DbPool.Query(ctx, sql)
	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	defer rows.Close()

	var heroes []*pbRepo.Hero
	for rows.Next() {
		h := &pbRepo.Hero{}

		err := rows.Scan(
			&h.Id,
			&h.Name,
			&h.BaseHealth,
			&h.BasePower,
			&h.BaseStrength,
			&h.BaseSpellpower,
			&h.HealthPerLevel,
			&h.PowerPerLevel,
			&h.StrengthPerLevel,
			&h.SpellpowerPerLevel,
		)

		if err != nil {
			s.Logger.Error(err)
			return nil, err
		}

		heroes = append(heroes, h)
	}

	return &pb.ListHeroesResponse{Heroes: heroes}, nil
}

func (s *CharacterService) ListCharacters(ctx context.Context, req *pb.ListCharactersRequest) (*pb.ListCharactersResponse, error) {
	sql, _, err := sq.
		Select("c.*, h.*").
		From("users as u").
		InnerJoin("characters as c on c.user_id = u.id").
		InnerJoin("heroes as h on h.id = c.hero_id").
		Where("u.username LIKE $1").ToSql()

	rows, err := s.DbPool.Query(ctx, sql, req.GetUsername())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var characters []*pbRepo.Character
	for rows.Next() {
		c := &pbRepo.Character{}
		h := &pbRepo.Hero{}

		err := rows.Scan(
			&c.CharacterId,
			&c.UserId,
			&c.CharacterName,
			&c.HeroId,
			&c.Level,
			&c.Experience,
			&c.Gold,
			&c.Tokens,
			&c.PointsHealth,
			&c.PointsPower,
			&c.PointsStrength,
			&c.PointsSpellpower,
			&c.UnspentPoints,
			&h.Id,
			&h.Name,
			&h.BaseHealth,
			&h.BasePower,
			&h.BaseStrength,
			&h.BaseSpellpower,
			&h.HealthPerLevel,
			&h.PowerPerLevel,
			&h.StrengthPerLevel,
			&h.SpellpowerPerLevel,
		)

		if err != nil {
			s.Logger.Error("failed to scan", "reason", err)
			return nil, err
		}

		c.Hero = h

		characters = append(characters, c)
	}

	return &pb.ListCharactersResponse{Characters: characters}, nil
}

func (s *CharacterService) SetSelectedCharacter(ctx context.Context, req *pb.SetSelectedCharacterRequest) (*pb.SetSelectedCharacterResponse, error) {
	sql, _, err := sq.
		Update("users").
		Set("selected_character_id", req.GetCharacterId()).
		Where("id = $2").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	_, err = s.DbPool.Exec(ctx, sql, req.GetCharacterId(), req.GetUserId())
	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	return &pb.SetSelectedCharacterResponse{Success: true, Message: "character selected."}, nil
}

func (s *CharacterService) GetLastSelectedCharacter(ctx context.Context, req *pb.GetLastSelectedCharacterRequest) (*pb.GetLastSelectedCharacterResponse, error) {
	sql, _, err := sq.
		Select("c.*, h.*").
		From("users AS u").
		InnerJoin("characters AS c ON c.character_id = u.selected_character_id").
		InnerJoin("heroes AS h ON h.id = c.hero_id").
		Where("lower(u.username) LIKE lower($1)").
		Limit(1).
		ToSql()

	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	rows, err := s.DbPool.Query(ctx, sql, req.GetUsername())
	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}
	defer rows.Close()

	var characters []*pbRepo.Character
	for rows.Next() {
		c := &pbRepo.Character{}
		h := &pbRepo.Hero{}

		err := rows.Scan(
			&c.CharacterId,
			&c.UserId,
			&c.CharacterName,
			&c.HeroId,
			&c.Level,
			&c.Experience,
			&c.Gold,
			&c.Tokens,
			&c.PointsHealth,
			&c.PointsPower,
			&c.PointsStrength,
			&c.PointsSpellpower,
			&c.UnspentPoints,
			&h.Id,
			&h.Name,
			&h.BaseHealth,
			&h.BasePower,
			&h.BaseStrength,
			&h.BaseSpellpower,
			&h.HealthPerLevel,
			&h.PowerPerLevel,
			&h.StrengthPerLevel,
			&h.SpellpowerPerLevel,
		)

		if err != nil {
			s.Logger.Error(err)
			return nil, err
		}

		c.Hero = h
		characters = append(characters, c)
	}

	if len(characters) == 0 {
		// select first character
		sql, _, err = sq.
			Select("c.*, h.*").
			From("users AS u").
			InnerJoin("characters AS c ON c.user_id = u.id").
			InnerJoin("heroes AS h ON h.id = c.hero_id").
			Where("lower(u.username) LIKE lower($1)").
			Limit(1).
			ToSql()
		if err != nil {
			s.Logger.Error(err)
			return nil, err
		}
		row := s.DbPool.QueryRow(ctx, sql, req.GetUsername())

		h := &pbRepo.Hero{}
		c := &pbRepo.Character{}

		err := row.Scan(
			&c.CharacterId,
			&c.UserId,
			&c.CharacterName,
			&c.HeroId,
			&c.Level,
			&c.Experience,
			&c.Gold,
			&c.Tokens,
			&c.PointsHealth,
			&c.PointsPower,
			&c.PointsStrength,
			&c.PointsSpellpower,
			&c.UnspentPoints,
			&h.Id,
			&h.Name,
			&h.BaseHealth,
			&h.BasePower,
			&h.BaseStrength,
			&h.BaseSpellpower,
			&h.HealthPerLevel,
			&h.PowerPerLevel,
			&h.StrengthPerLevel,
			&h.SpellpowerPerLevel,
		)
		if err != nil {
			s.Logger.Error(err)
			return nil, err
		}

		c.Hero = h

		characters = append(characters, c)
	}

	return &pb.GetLastSelectedCharacterResponse{
		Character: characters[0],
	}, nil
}

func (s *CharacterService) GetCharacter(ctx context.Context, req *pb.GetCharacterRequest) (*pb.GetCharacterResponse, error) {
	sql, _, err := sq.
		Select("c.*, h.*").
		From("characters AS c").
		InnerJoin("heroes AS h ON h.id = c.hero_id").
		Where("c.character_id = $1").
		Limit(1).
		ToSql()

	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	row := s.DbPool.QueryRow(ctx, sql, req.GetCharacterId())

	h := &pbRepo.Hero{}
	c := &pbRepo.Character{}

	err = row.Scan(
		&c.CharacterId,
		&c.UserId,
		&c.CharacterName,
		&c.HeroId,
		&c.Level,
		&c.Experience,
		&c.Gold,
		&c.Tokens,
		&c.PointsHealth,
		&c.PointsPower,
		&c.PointsStrength,
		&c.PointsSpellpower,
		&c.UnspentPoints,
		&h.Id,
		&h.Name,
		&h.BaseHealth,
		&h.BasePower,
		&h.BaseStrength,
		&h.BaseSpellpower,
		&h.HealthPerLevel,
		&h.PowerPerLevel,
		&h.StrengthPerLevel,
		&h.SpellpowerPerLevel,
	)

	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	c.Hero = h

	return &pb.GetCharacterResponse{Character: c}, nil
}

func (s *CharacterService) CreateCharacter(ctx context.Context, req *pb.CreateCharacterRequest) (*pb.CreateCharacterResponse, error) {
	capitalizedName := fmt.Sprint(strings.ToUpper(req.Name[:1]), req.Name[1:])

	sql, args, err := sq.
		Insert("characters").
		Columns("user_id", "character_name", "hero_id", "level", "exp", "gold", "tokens", "points_health", "points_power", "points_strength", "points_spellpower").
		Values(req.GetUserId(), capitalizedName, req.GetHero().GetId(), 1, 1, 50, 0, 0, 0, 0, 0).
		Suffix("RETURNING character_id").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	row := s.DbPool.QueryRow(ctx, sql, args...)

	var characterId int64
	err = row.Scan(&characterId)
	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	sql, _, err = sq.
		Select("c.*, h.*").
		From("characters AS c").
		InnerJoin("heroes AS h ON h.id = c.hero_id").
		Where("c.character_id = $1").
		ToSql()

	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	row = s.DbPool.QueryRow(ctx, sql, characterId)
	h := &pbRepo.Hero{}
	c := &pbRepo.Character{}
	err = row.Scan(
		&c.CharacterId,
		&c.UserId,
		&c.CharacterName,
		&c.HeroId,
		&c.Level,
		&c.Experience,
		&c.Gold,
		&c.Tokens,
		&c.PointsHealth,
		&c.PointsPower,
		&c.PointsStrength,
		&c.PointsSpellpower,
		&c.UnspentPoints,
		&h.Id,
		&h.Name,
		&h.BaseHealth,
		&h.BasePower,
		&h.BaseStrength,
		&h.BaseSpellpower,
		&h.HealthPerLevel,
		&h.PowerPerLevel,
		&h.StrengthPerLevel,
		&h.SpellpowerPerLevel,
	)

	c.Hero = h

	// create empty equipped slots for character
	// get equipment slots
	sql, _, err = sq.Select("*").From("equip_slots").ToSql()

	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	rows, err := s.DbPool.Query(ctx, sql)
	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}
	defer rows.Close()

	var equipSlots []*pbRepo.EquipSlot
	for rows.Next() {
		es := &pbRepo.EquipSlot{}

		err := rows.Scan(
			&es.Id,
			&es.Code,
			&es.Name,
		)

		if err != nil {
			s.Logger.Error(err)
			return nil, err
		}

		equipSlots = append(equipSlots, es)
	}

	for _, equipSlot := range equipSlots {
		_, err = s.DbPool.Exec(ctx, "INSERT INTO character_equipments (character_id, equip_slot_id) VALUES ($1, $2)", characterId, equipSlot.Id)
		if err != nil {
			s.Logger.Error(err)
			return nil, err
		}
	}

	return &pb.CreateCharacterResponse{Success: true, Message: "character created.", Character: c}, nil
}

func (s *CharacterService) DeleteCharacter(ctx context.Context, req *pb.DeleteCharacterRequest) (*pb.DeleteCharacterResponse, error) {
	sql, _, err := sq.Delete("*").From("characters").Where("character_id = $1").ToSql()
	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	_, err = s.DbPool.Exec(ctx, sql, req.GetCharacterId())
	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	return &pb.DeleteCharacterResponse{Success: true, Message: "character deleted."}, nil
}

func (s *CharacterService) GetEquippedItems(ctx context.Context, req *pb.GetEquippedItemsRequest) (*pb.GetEquippedItemsResponse, error) {
	sql, _, err := sq.
		Select("ce.*, ii.*, i.*, t.*, es.*").
		From("character_equipments AS ce").
		InnerJoin("equip_slots AS es ON es.id = ce.equip_slot_id").
		InnerJoin("item_instances AS ii ON ii.item_instance_id = ce.item_instance_id").
		InnerJoin("items AS i ON i.id = ii.item_id").
		InnerJoin("item_templates AS t ON t.id = i.item_template_id").
		Where("ce.character_id = $1").
		ToSql()

	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	rows, err := s.DbPool.Query(ctx, sql, req.GetCharacterId())
	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}
	defer rows.Close()

	var itemInstances []*pbRepo.ItemInstance
	for rows.Next() {
		var foo *any
		es := &pbRepo.EquipSlot{}
		ii := &pbRepo.ItemInstance{}
		ce := &pbRepo.CharacterEquipment{}
		t := &pbRepo.ItemTemplate{}
		i := &pbRepo.Item{}

		err := rows.Scan(
			&ce.CharacterId,
			&ce.EquipSlotId,
			&ce.ItemInstanceId,
			&ii.ItemInstanceId,
			&ii.ItemId,
			&ii.OwnerCharacterId,
			&foo,
			&foo,
			&i.Id,
			&i.Prefix,
			&i.Suffix,
			&i.ItemTemplateId,
			&i.Health,
			&i.Power,
			&i.Strength,
			&i.Spellpower,
			&i.BonusDamage,
			&i.BonusArmor,
			&t.Id,
			&t.Name,
			&t.ItemTypeId,
			&t.EquipSlotId,
			&t.Description,
			&t.GoldPrice,
			&t.BuyableWithTokens,
			&t.TokenPrice,
			&foo,
			&es.Id,
			&es.Code,
			&es.Name,
		)

		t.EquipSlot = es
		i.ItemTemplate = t
		ii.Item = i
		ii.OwnerCharacterId = req.CharacterId

		if err != nil {
			s.Logger.Error(err)
			return nil, err
		}

		itemInstances = append(itemInstances, ii)
	}

	return &pb.GetEquippedItemsResponse{ItemInstances: itemInstances}, nil
}

func (s *CharacterService) EquipItem(ctx context.Context, req *pb.EquipItemRequest) (*pb.EquipItemResponse, error) {
	sql, _, err := sq.Update("character_equipments").
		Set("item_instance_id", req.ItemInstanceId).
		Where("character_id = $2").
		Where("equip_slot_id = $3").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	_, err = s.DbPool.Exec(ctx, sql, req.ItemInstanceId, req.CharacterId, req.EquipSlotId)
	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	return &pb.EquipItemResponse{Success: true}, nil
}

func (s *CharacterService) UnequipItem(ctx context.Context, req *pb.UnequipItemRequest) (*pb.UnequipItemResponse, error) {
	sql, _, err := sq.Update("character_equipments").
		Set("item_instance_id", nil).
		Where("character_id = $2").
		Where("equip_slot_id = $3").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	_, err = s.DbPool.Exec(ctx, sql, nil, req.CharacterId, req.EquipSlotId)
	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	return &pb.UnequipItemResponse{Success: true}, nil
}

func (s *CharacterService) ToggleEquip(ctx context.Context, req *pb.ToggleEquipRequest) (*pb.ToggleEquipResponse, error) {
	if req.ShouldEquip {
		req := &pb.EquipItemRequest{CharacterId: req.CharacterId, ItemInstanceId: req.ItemInstanceId, EquipSlotId: req.EquipSlotId}
		_, err := s.EquipItem(ctx, req)
		if err != nil {
			s.Logger.Error(err)
			return nil, err
		}
		return &pb.ToggleEquipResponse{Success: true, Message: "item equipped."}, nil
	} else {
		req := &pb.UnequipItemRequest{CharacterId: req.CharacterId, EquipSlotId: req.EquipSlotId}
		_, err := s.UnequipItem(ctx, req)
		if err != nil {
			s.Logger.Error(err)
			return nil, err
		}
		return &pb.ToggleEquipResponse{Success: true, Message: "item unequipped."}, nil
	}
}

func (s *CharacterService) GetEquipSlots(ctx context.Context, req *pb.GetEquipSlotsRequest) (*pb.GetEquipSlotsResponse, error) {
	sql, _, err := sq.
		Select("*").
		From("equip_slots").
		ToSql()

	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	rows, err := s.DbPool.Query(ctx, sql)
	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}
	defer rows.Close()

	var equipSlots []*pbRepo.EquipSlot
	for rows.Next() {
		es := &pbRepo.EquipSlot{}

		err := rows.Scan(
			&es.Id,
			&es.Code,
			&es.Name,
		)

		if err != nil {
			s.Logger.Error(err)
			return nil, err
		}

		equipSlots = append(equipSlots, es)
	}

	return &pb.GetEquipSlotsResponse{Slots: equipSlots}, nil
}

func (s *CharacterService) SaveCombatWinProgress(ctx context.Context, req *pb.SaveCombatWinProgressRequest) (*pb.SaveCombatWinProgressResponse, error) {
	char := req.Character
	char = functions.AddXp(req.Character, req.Npc.NpcTemplate.BaseXpReward)

	// write char to db
	sql, args, err := sq.
		Update("characters").
		Set("exp", char.Experience).
		Set("level", char.Level).
		Set("points_health", char.PointsHealth).
		Set("points_power", char.PointsPower).
		Set("points_strength", char.PointsStrength).
		Set("points_spellpower", char.PointsSpellpower).
		Where(sq.Eq{"character_id": char.CharacterId}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	_, err = s.DbPool.Exec(ctx, sql, args...)
	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	sql, _, err = sq.
		Select("c.*, h.*").
		From("characters AS c").
		InnerJoin("heroes AS h ON h.id = c.hero_id").
		Where("c.character_id = $1").
		Limit(1).
		ToSql()

	row := s.DbPool.QueryRow(ctx, sql, char.CharacterId)

	char.Hero = &pbRepo.Hero{}
	err = row.Scan(
		&char.CharacterId,
		&char.UserId,
		&char.CharacterName,
		&char.HeroId,
		&char.Level,
		&char.Experience,
		&char.Gold,
		&char.Tokens,
		&char.PointsHealth,
		&char.PointsPower,
		&char.PointsStrength,
		&char.PointsSpellpower,
		&char.UnspentPoints,
		&char.Hero.Id,
		&char.Hero.Name,
		&char.Hero.BaseHealth,
		&char.Hero.BasePower,
		&char.Hero.BaseStrength,
		&char.Hero.BaseSpellpower,
		&char.Hero.HealthPerLevel,
		&char.Hero.PowerPerLevel,
		&char.Hero.StrengthPerLevel,
		&char.Hero.SpellpowerPerLevel,
	)

	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	return &pb.SaveCombatWinProgressResponse{
		Success:   true,
		Message:   "character update saved",
		Character: char,
	}, nil
}

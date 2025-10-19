package server

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	sq "github.com/Masterminds/squirrel"

	pb "github.com/komadiina/spelltext/proto/armory"
	pbRepo "github.com/komadiina/spelltext/proto/repo"
	"github.com/komadiina/spelltext/server/character/config"
	"github.com/komadiina/spelltext/utils/singleton/logging"
)

type CharacterService struct {
	pb.UnimplementedCharacterServer
	DbPool *pgxpool.Pool
	Config *config.Config
	Logger *logging.Logger
}

func tryConnect(s *CharacterService, context context.Context, conninfo string, backoff time.Duration, maxRetries int, boFormula func(time.Duration) time.Duration) (pgx.Conn, error) {
	attempt := 1
	for {
		conn, err := pgx.Connect(context, conninfo)

		if err != nil && attempt >= maxRetries {
			// conn not established, max retries exceeded
			s.Logger.Fatal(err)
		} else if err == nil && attempt < maxRetries {
			// conn established within maxRetries
			s.Logger.Info("pgpool connection established")
			return *conn, nil
		} else if err != nil && attempt < maxRetries {
			// conn not established, backoff
			s.Logger.Warn("failed to establish database connection, backing off...", "reason", err, "backoff_seconds", backoff.Seconds())
			time.Sleep(backoff)
			backoff = boFormula(backoff)
			attempt++
		}
	}
}

func (s *CharacterService) GetConn(ctx context.Context) *pgx.Conn {
	conninfo := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		s.Config.PgUser,
		s.Config.PgPass,
		s.Config.PgHost,
		s.Config.PgPort,
		s.Config.PgDbName,
		s.Config.PgSSLMode,
	)

	backoff := time.Second * 5 // secs
	time.Sleep(backoff)

	conn, err := tryConnect(s, ctx, conninfo, backoff, 5, func(backoff time.Duration) time.Duration {
		backoff = backoff + time.Second*5
		return backoff
	})

	if err != nil {
		return nil
	} else {
		s.Logger.Error(err)
	}

	return &conn
}

func (s *CharacterService) ListHeroes(ctx context.Context, req *pb.ListHeroesRequest) (*pb.ListHeroesResponse, error) {
	s.Logger.Warn("unimplemented method called", "name", "ListHeroes")

	return nil, nil
}

func (s *CharacterService) ListCharacters(ctx context.Context, req *pb.ListCharactersRequest) (*pb.ListCharactersResponse, error) {
	cte, _, err := sq.Select("u.id AS id, u.username AS username").
		From("users AS u").
		Where("u.username LIKE $1").
		ToSql()

	if err != nil {
		s.Logger.Error("failed to build cte", "err", err)
		return nil, nil
	}

	query, _, err := sq.Select("c.*, h.*").
		From("characters AS c").
		InnerJoin("u_filt ON u_filt.username LIKE $2").
		InnerJoin("heroes AS h ON h.id = c.character_id").ToSql()

	if err != nil {
		s.Logger.Error("failed to build query", "err", err)
		return nil, nil
	}

	sql := fmt.Sprintf("WITH u_filt AS (%s) %s", cte, query)
	rows, err := s.DbPool.Query(ctx, sql, req.GetUsername(), req.GetUsername())

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

		if err != nil {
			s.Logger.Error(err)
			return nil, err
		}

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
	sql, _, err := sq.
		Insert("characters").
		Columns("user_id", "character_name", "hero_id", "level", "experience", "gold", "tokens", "points_health", "points_power", "points_strength", "points_spellpower").
		Values(req.GetUserId(), req.GetName(), req.GetHero().GetId(), 1, 0, 0, 0, 0, 0, 0, 0).
		Suffix("RETURNING character_id").
		ToSql()

	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	row := s.DbPool.QueryRow(ctx, sql)

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

	var items []*pbRepo.Item
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
			&ii.CreatedAt,
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

		if err != nil {
			s.Logger.Error(err)
			return nil, err
		}

		items = append(items, i)
	}

	return &pb.GetEquippedItemsResponse{Items: items}, nil
}

func (s *CharacterService) EquipItem(ctx context.Context, req *pb.EquipItemRequest) (*pb.EquipItemResponse, error) {
	sql, _, err := sq.Update("character_equipment").
		Set("item_instance_id", req.ItemInstanceId).
		Where("character_id = $2").
		Where("equip_slot_id = $3").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	s.Logger.Info(sql)

	_, err = s.DbPool.Exec(ctx, sql, req.ItemInstanceId, req.CharacterId, req.EquipSlotId)
	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	return &pb.EquipItemResponse{Success: true}, nil
}

func (s *CharacterService) UnequipItem(ctx context.Context, req *pb.UnequipItemRequest) (*pb.UnequipItemResponse, error) {
	sql, _, err := sq.Update("character_equipment").
		Set("item_instance_id", nil).
		Where("character_id = $2").
		Where("equip_slot_id = $3").
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

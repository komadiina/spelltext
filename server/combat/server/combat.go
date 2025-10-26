package server

import (
	"context"
	"errors"
	"math/rand"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	pbChar "github.com/komadiina/spelltext/proto/char"
	pb "github.com/komadiina/spelltext/proto/combat"
	pbHealth "github.com/komadiina/spelltext/proto/health"
	pbRepo "github.com/komadiina/spelltext/proto/repo"
	"github.com/komadiina/spelltext/server/combat/config"
	health "github.com/komadiina/spelltext/utils"
	"github.com/komadiina/spelltext/utils/singleton/logging"
	"google.golang.org/grpc"
)

type Clients struct {
	Character pbChar.CharacterClient
}

type Connections struct {
	Character *grpc.ClientConn
}

type CombatService struct {
	pb.UnimplementedCombatServer
	health.HealthCheckable
	DbPool      *pgxpool.Pool
	Config      *config.Config
	Logger      *logging.Logger
	Connections *Connections
	Clients     *Clients
}

func (s *CombatService) Check(ctx context.Context, req *pbHealth.HealthCheckRequest) (*pbHealth.HealthCheckResponse, error) {
	return &pbHealth.HealthCheckResponse{Status: pbHealth.HealthCheckResponse_SERVING}, nil
}

func (s *CombatService) ListNpcs(ctx context.Context, req *pb.ListNpcsRequest) (*pb.ListNpcsResponse, error) {
	sql, args, err := sq.
		Select("npc.*, npc_t.*").
		From("npcs as npc").
		InnerJoin("npc_templates as npc_t on npc_t.id = npc.template_id").
		InnerJoin("characters ch on ch.character_id = $1", req.CharacterId).
		Where("npc_t.min_level <= ch.level AND npc_t.max_level >= ch.level").
		ToSql()

	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	rows, err := s.DbPool.Query(ctx, sql, args...)
	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}
	defer rows.Close()

	var nli []*pb.NpcListItem

	var npcs []*pbRepo.Npc
	for rows.Next() {
		npc := &pbRepo.Npc{}
		npcTempl := &pbRepo.NpcTemplate{}

		err := rows.Scan(
			&npc.Id,
			&npc.Prefix,
			&npc.Suffix,
			&npc.TemplateId,
			&npc.HealthMultiplier,
			&npc.DamageMultiplier,
			&npcTempl.Id,
			&npcTempl.Name,
			&npcTempl.Description,
			&npcTempl.MinLevel,
			&npcTempl.MaxLevel,
			&npcTempl.HealthPoints,
			&npcTempl.BaseDamage,
			&npcTempl.BaseXpReward,
			&npcTempl.GoldReward,
		)

		if err != nil {
			s.Logger.Error(err)
			return nil, err
		}

		npc.NpcTemplate = npcTempl
		npc.Level = rand.Uint32()%(npcTempl.MaxLevel-npcTempl.MinLevel+1) + npcTempl.MinLevel
		npcs = append(npcs, npc)
	}

	var items []*pbRepo.Item
	// query for an item in db by npc_template_id == returns 1 item associated with npc loot table
	sql, _, err = sq.
		Select("nlt.*, i.*, templ.*, it.*, es.*").
		From("npc_loot_tables as nlt").
		InnerJoin("items as i on i.id = nlt.item_id").
		InnerJoin("item_templates as templ on templ.id = i.item_template_id").
		InnerJoin("equip_slots as es on es.id = templ.equip_slot_id").
		InnerJoin("item_types as it on it.id = templ.item_type_id").
		Where("nlt.npc_template_id = $1").
		ToSql()

	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	for _, npc := range npcs {
		clear(items)

		rows, err := s.DbPool.Query(ctx, sql, npc.TemplateId)
		if err != nil {
			s.Logger.Error(err)
			return nil, err
		}

		items = make([]*pbRepo.Item, 0)
		for rows.Next() {
			var foo *any
			nlt := &pbRepo.NpcLootTable{}
			i := &pbRepo.Item{}
			templ := &pbRepo.ItemTemplate{}
			it := &pbRepo.ItemType{}
			es := &pbRepo.EquipSlot{}

			err := rows.Scan(
				&nlt.NpcTemplateId,
				&nlt.ItemId,
				&nlt.Chance,
				&nlt.Quantity,
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
				&templ.Id,
				&templ.Name,
				&templ.ItemTypeId,
				&templ.EquipSlotId,
				&templ.Description,
				&templ.GoldPrice,
				&templ.BuyableWithTokens,
				&templ.TokenPrice,
				&foo,
				&it.Id,
				&it.Code,
				&it.Name,
				&es.Id,
				&es.Code,
				&es.Name,
			)
			if err != nil {
				s.Logger.Error(err)
				return nil, err
			}

			i.ItemTemplate = templ
			templ.EquipSlot = es
			templ.ItemType = it

			items = append(items, i)
		}

		nli = append(nli, &pb.NpcListItem{Npc: npc, Drops: items})
		rows.Close()
	}

	return &pb.ListNpcsResponse{Npcs: nli}, nil
}

func (s *CombatService) SubmitLoss(ctx context.Context, req *pb.SubmitLossRequest) (*pb.SubmitLossResponse, error) {
	// nothing atm, TODO: add combat_history table (character_id, npc_id, player_lost-bool)
	return nil, nil
}

func (s *CombatService) SubmitWin(ctx context.Context, req *pb.SubmitWinRequest) (*pb.SubmitWinResponse, error) {
	s.Logger.Infof("req: %+v", req)
	sql, _, err := sq.
		Select("c.*, h.*").
		From("characters as c").
		InnerJoin("heroes as h on h.id = c.hero_id").
		Where("c.character_id = $1").
		Limit(1).
		ToSql()


	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	row := s.DbPool.QueryRow(ctx, sql, req.GetCharacterId())

	char := &pbRepo.Character{}
	hero := &pbRepo.Hero{}
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
		&hero.Id,
		&hero.Name,
		&hero.BaseHealth,
		&hero.BasePower,
		&hero.BaseStrength,
		&hero.BaseSpellpower,
		&hero.HealthPerLevel,
		&hero.PowerPerLevel,
		&hero.StrengthPerLevel,
		&hero.SpellpowerPerLevel,
	)

	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	char.Hero = hero

	sql, _, err = sq.
		Select("npc.*, npc_t.*").
		From("npcs as npc").
		InnerJoin("npc_templates as npc_t on npc_t.id = npc.template_id").
		Where("npc.id = $1").
		Limit(1).
		ToSql()

	row = s.DbPool.QueryRow(ctx, sql, req.GetNpcId())
	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	npc := &pbRepo.Npc{}
	npcTempl := &pbRepo.NpcTemplate{}
	err = row.Scan(
		&npc.Id,
		&npc.Prefix,
		&npc.Suffix,
		&npc.TemplateId,
		&npc.HealthMultiplier,
		&npc.DamageMultiplier,
		&npcTempl.Id,
		&npcTempl.Name,
		&npcTempl.Description,
		&npcTempl.MinLevel,
		&npcTempl.MaxLevel,
		&npcTempl.HealthPoints,
		&npcTempl.BaseDamage,
		&npcTempl.BaseXpReward,
		&npcTempl.GoldReward,
	)
	npc.NpcTemplate = npcTempl

	var items []*pbRepo.Item
	// query for an item in db by npc_template_id == returns 1 item associated with npc loot table
	sql, _, err = sq.
		Select("nlt.*, i.*, templ.*, it.*, es.*").
		From("npc_loot_tables as nlt").
		InnerJoin("items as i on i.id = nlt.item_id").
		InnerJoin("item_templates as templ on templ.id = i.item_template_id").
		InnerJoin("equip_slots as es on es.id = templ.equip_slot_id").
		InnerJoin("item_types as it on it.id = templ.item_type_id").
		Where("nlt.npc_template_id = $1").
		ToSql()

	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	rows, err := s.DbPool.Query(ctx, sql, npc.TemplateId)
	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	for rows.Next() {
		var foo *any
		nlt := &pbRepo.NpcLootTable{}
		i := &pbRepo.Item{}
		templ := &pbRepo.ItemTemplate{}
		it := &pbRepo.ItemType{}
		es := &pbRepo.EquipSlot{}

		err := rows.Scan(
			&nlt.NpcTemplateId,
			&nlt.ItemId,
			&nlt.Chance,
			&nlt.Quantity,
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
			&templ.Id,
			&templ.Name,
			&templ.ItemTypeId,
			&templ.EquipSlotId,
			&templ.Description,
			&templ.GoldPrice,
			&templ.BuyableWithTokens,
			&templ.TokenPrice,
			&foo,
			&it.Id,
			&it.Code,
			&it.Name,
			&es.Id,
			&es.Code,
			&es.Name,
		)
		if err != nil {
			s.Logger.Error(err)
			return nil, err
		}

		i.ItemTemplate = templ
		templ.EquipSlot = es
		templ.ItemType = it

		items = append(items, i)
	}

	resp, err := s.Clients.Character.SaveCombatWinProgress(
		ctx,
		&pbChar.SaveCombatWinProgressRequest{
			Character: char,
			Npc:       npc,
		},
	)

	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	if !resp.Success {
		return nil, errors.New(resp.Message)
	}

	// nothing else atm, TODO: add combat_history table (character_id, npc_id, player_lost-bool)
	return &pb.SubmitWinResponse{
		LevelUp:      resp.Character.Level != char.Level,
		NewCharacter: resp.Character,
		ItemReward:   items,
		GoldReward:   npc.NpcTemplate.GoldReward,
		XpReward:     npc.NpcTemplate.BaseXpReward,
	}, nil
}

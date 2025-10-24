package server

import (
	"context"
	"math/rand"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	pb "github.com/komadiina/spelltext/proto/combat"
	pbHealth "github.com/komadiina/spelltext/proto/health"
	pbRepo "github.com/komadiina/spelltext/proto/repo"
	"github.com/komadiina/spelltext/server/combat/config"
	health "github.com/komadiina/spelltext/utils"
	"github.com/komadiina/spelltext/utils/singleton/logging"
)

type CombatService struct {
	pb.UnimplementedCombatServer
	health.HealthCheckable
	DbPool *pgxpool.Pool
	Config *config.Config
	Logger *logging.Logger
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

			nli = append(nli, &pb.NpcListItem{Npc: npc, Drops: items})

			items = append(items, i)
		}

		rows.Close()
	}

	return &pb.ListNpcsResponse{Npcs: nli}, nil
}

func (s *CombatService) SubmitLoss(ctx context.Context, req *pb.SubmitLossRequest) (*pb.SubmitLossResponse, error) {
	// nothing atm, TODO: add combat_history table (character_id, npc_id, player_lost-bool)
	return nil, nil
}

func (s *CombatService) SubmitWin(ctx context.Context, req *pb.SubmitWinRequest) (*pb.SubmitWinResponse, error) {
	// gather items
	// sql, args, err := sq.Select("*").
	// 	From("").
	// ToSql()

	return nil, nil
}

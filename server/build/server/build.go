package server

import (
	"context"
	"fmt"
	"slices"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	pb "github.com/komadiina/spelltext/proto/build"
	pbCharacter "github.com/komadiina/spelltext/proto/char"
	pbHealth "github.com/komadiina/spelltext/proto/health"
	pbRepo "github.com/komadiina/spelltext/proto/repo"
	"github.com/komadiina/spelltext/server/build/config"
	generics "github.com/komadiina/spelltext/utils"
	health "github.com/komadiina/spelltext/utils"
	"github.com/komadiina/spelltext/utils/singleton/logging"
	"google.golang.org/grpc"
)

type Clients struct {
	Character pbCharacter.CharacterClient
}

type Connections struct {
	Character *grpc.ClientConn
}

type BuildService struct {
	pb.UnimplementedBuildServer
	health.HealthCheckable
	DbPool      *pgxpool.Pool
	Config      *config.Config
	Logger      *logging.Logger
	Clients     *Clients
	Connections *Connections
}

func (s *BuildService) Check(ctx context.Context, req *pbHealth.HealthCheckRequest) (*pbHealth.HealthCheckResponse, error) {
	return &pbHealth.HealthCheckResponse{Status: pbHealth.HealthCheckResponse_SERVING}, nil
}

func (s *BuildService) ListAbilities(ctx context.Context, req *pb.ListAbilitiesRequest) (*pb.ListAbilitiesResponse, error) {
	// get upgraded abilities (if any)
	sql, _, err := sq.
		Select("pat.*, a.*").
		From("player_ability_tree as pat").
		InnerJoin("abilities as a on a.id = pat.ability_id").
		Where("character_id = $1").
		ToSql()

	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	rows, err := s.DbPool.Query(ctx, sql, req.Character.CharacterId)
	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}
	defer rows.Close()

	var upgraded []*pbRepo.PlayerAbilityTree
	for rows.Next() {
		at := &pbRepo.PlayerAbilityTree{}
		a := &pbRepo.Ability{}
		err := rows.Scan(
			&at.CharacterId,
			&at.AbilityId,
			&at.Level,
			&a.Id,
			&a.Name,
			&a.Description,
			&a.Type,
			&a.TalentPointCost,
			&a.PowerCost,
			&a.BaseDamage,
			&a.StrengthMultiplier,
			&a.SpellpowerMultiplier,
			&a.StMultPerlevel,
			&a.SpMultPerlevel,
			&a.MinLevel,
		)
		if err != nil {
			s.Logger.Error(err)
			return nil, err
		}

		at.Ability = a
		upgraded = append(upgraded, at)
	}

	mapped := generics.Map(upgraded, func(a *pbRepo.PlayerAbilityTree) uint64 { return a.AbilityId })

	sql, _, err = sq.Select("*").From("abilities").ToSql()
	if err != nil {
		s.Logger.Error("failed to build query", "err", err)
		return nil, err
	}

	rows, err = s.DbPool.Query(ctx, sql)
	if err != nil {
		s.Logger.Error("failed to run query", "err", err)
		return nil, err
	}
	defer rows.Close()

	var available []*pbRepo.Ability
	var locked []*pbRepo.Ability
	for rows.Next() {
		a := &pbRepo.Ability{}
		err := rows.Scan(
			&a.Id,
			&a.Name,
			&a.Description,
			&a.Type,
			&a.TalentPointCost,
			&a.PowerCost,
			&a.BaseDamage,
			&a.StrengthMultiplier,
			&a.SpellpowerMultiplier,
			&a.StMultPerlevel,
			&a.SpMultPerlevel,
			&a.MinLevel,
		)

		if err != nil {
			s.Logger.Error("failed to scan", "err", err)
			return nil, err
		}

		if a.MinLevel > req.Character.Level {
			locked = append(locked, a)
		} else if !slices.Contains(mapped, a.Id) {
			available = append(available, a)
		}
	}

	return &pb.ListAbilitiesResponse{Available: available, Locked: locked, Upgraded: upgraded}, nil
}

func (s *BuildService) UpgradeAbility(ctx context.Context, req *pb.UpgradeAbilityRequest) (*pb.UpgradeAbilityResponse, error) {
	resp, err := s.Clients.Character.GetCharacter(ctx, &pbCharacter.GetCharacterRequest{CharacterId: req.CharacterId})
	if err != nil {
		s.Logger.Error("failed to get character", "err", err)
		return nil, err
	}

	if resp.Character.UnspentPoints <= 0 {
		s.Logger.Warnf("unable to level up ability (%d), not enough unspent points: %d, character_id=%d", req.AbilityId, resp.Character.UnspentPoints, req.CharacterId)
		return nil, fmt.Errorf("not enough points")
	}

	tx, err := s.DbPool.BeginTx(ctx, pgx.TxOptions{})
	defer tx.Rollback(ctx)

	b := &pgx.Batch{}
	if req.NewAbility {
		b.Queue("INSERT INTO player_ability_tree (character_id, ability_id, level) VALUES ($1, $2, 1)", req.CharacterId, req.AbilityId)
	} else {
		b.Queue("UPDATE player_ability_tree SET level = level + 1 WHERE character_id = $1 AND ability_id = $2", req.CharacterId, req.AbilityId)
	}
	b.Queue("UPDATE characters SET unspent_points = (CASE WHEN unspent_points = 0 THEN 1 ELSE unspent_points END) - 1 WHERE character_id = $1", req.CharacterId)

	res := s.DbPool.SendBatch(ctx, b)
	if err = res.Close(); err != nil {
		s.Logger.Error(err)
		return nil, nil
	}

	tx.Commit(ctx)

	return &pb.UpgradeAbilityResponse{Success: true}, nil
}

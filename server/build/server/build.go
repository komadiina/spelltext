package server

import (
	"context"
	"fmt"
	"slices"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	pb "github.com/komadiina/spelltext/proto/build"
	pbCharacter "github.com/komadiina/spelltext/proto/char"
	pbRepo "github.com/komadiina/spelltext/proto/repo"
	"github.com/komadiina/spelltext/server/build/config"
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
	DbPool      *pgxpool.Pool
	Config      *config.Config
	Logger      *logging.Logger
	Clients     *Clients
	Connections *Connections
}

func (s *BuildService) ListAbilities(ctx context.Context, req *pb.ListAbilitiesRequest) (*pb.ListAbilitiesResponse, error) {
	// get upgraded abilities (if any)
	sql, _, err := sq.Select("ability_id").From("player_ability_tree").Where("character_id = $1").ToSql()
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

	var upgradedAbilities []uint64
	for rows.Next() {
		var abilityId uint64
		err := rows.Scan(&abilityId)
		if err != nil {
			s.Logger.Error(err)
			return nil, err
		}

		upgradedAbilities = append(upgradedAbilities, abilityId)
	}

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
	var upgraded []*pbRepo.Ability
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
		} else if !slices.Contains(upgradedAbilities, a.Id) {
			available = append(available, a)
		} else {
			upgraded = append(upgraded, a)
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

	sql := "UPDATE character_ability_tree SET level = level + 1 WHERE character_id = $1 AND ability_id = $2"
	_, err = s.DbPool.Exec(ctx, sql, req.CharacterId, req.AbilityId)
	return &pb.UpgradeAbilityResponse{Success: true}, nil
}

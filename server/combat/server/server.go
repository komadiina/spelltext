package server

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	pb "github.com/komadiina/spelltext/proto/combat"
	"github.com/komadiina/spelltext/server/combat/config"
	"github.com/komadiina/spelltext/utils/singleton/logging"
)

type CombatService struct {
	pb.UnimplementedCombatServer
	DbPool *pgxpool.Pool
	Config *config.Config
	Logger *logging.Logger
}

func (s *CombatService) ListNpcs(ctx context.Context, req *pb.ListNpcsRequest) (*pb.ListNpcsResponse, error) {
	return nil, nil
}

func (s *CombatService) InitiateCombat(ctx context.Context, req *pb.InitiateCombatRequest) (*pb.InitiateCombatResponse, error) {
	return nil, nil
}

func (s *CombatService) ResolveCombat(ctx context.Context, req *pb.ResolveCombatRequest) (*pb.ResolveCombatResponse, error) {
	return nil, nil
}

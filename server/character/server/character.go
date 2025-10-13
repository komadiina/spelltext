package server

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	sq "github.com/Masterminds/squirrel"

	pb "github.com/komadiina/spelltext/proto/armory"
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

	query, _, err := sq.Select("c.character_id, c.character_name, h.name, c.level, c.gold, c.tokens").
		From("characters AS c").
		InnerJoin("u_filt ON u_filt.username LIKE $2").
		InnerJoin("heroes AS h ON h.id = c.character_id").ToSql()

	if err != nil {
		s.Logger.Error("failed to build query", "err", err)
		return nil, nil
	}

	sql := fmt.Sprintf("WITH u_filt AS (%s) %s", cte, query)
	rows, err := s.DbPool.Query(ctx, sql, req.GetUsername(), req.GetUsername())

	var characters []*pb.TCharacter
	for rows.Next() {
		c := &pb.TCharacter{}

		err := rows.Scan(
			&c.Id,
			&c.Name,
			&c.Class,
			&c.Level,
			&c.Gold,
			&c.Tokens,
		)

		if err != nil {
			s.Logger.Error("failed to scan", "reason", err)
			return nil, err
		}

		characters = append(characters, c)
	}

	return &pb.ListCharactersResponse{Characters: characters}, nil
}

func (s *CharacterService) CreateCharacter(ctx context.Context, req *pb.CreateCharacterRequest) (*pb.CreateCharacterResponse, error) {

	return nil, nil
}

func (s *CharacterService) DeleteCharacter(ctx context.Context, req *pb.DeleteCharacterRequest) (*pb.DeleteCharacterResponse, error) {
	s.Logger.Warn("unimplemented method called", "name", "DeleteCharacter")
	return nil, nil
}

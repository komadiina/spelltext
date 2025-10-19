package server

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	pbArmory "github.com/komadiina/spelltext/proto/armory"
	pb "github.com/komadiina/spelltext/proto/auth"
	pbRepo "github.com/komadiina/spelltext/proto/repo"
	"github.com/komadiina/spelltext/server/auth/config"
	"github.com/komadiina/spelltext/utils/singleton/logging"
	"google.golang.org/grpc"
)

type Clients struct {
	Armory pbArmory.CharacterClient
}

type Connections struct {
	Armory *grpc.ClientConn
}

type AuthService struct {
	pb.UnimplementedAuthServer
	DbPool      *pgxpool.Pool
	Config      *config.Config
	Logger      *logging.Logger
	Clients     *Clients
	Connections *Connections
}

func tryConnect(s *AuthService, context context.Context, conninfo string, backoff time.Duration, maxRetries int, boFormula func(time.Duration) time.Duration) (pgx.Conn, error) {
	try := 1
	for {
		conn, err := pgx.Connect(context, conninfo)

		if err != nil && try >= maxRetries {
			// conn not established, max retries exceeded
			s.Logger.Fatal(err)
		} else if err == nil && try < maxRetries {
			// conn established within maxRetries
			s.Logger.Info("pgpool connection established")
			return *conn, nil
		} else if err != nil && try < maxRetries {
			// conn not established, backoff
			s.Logger.Warn("failed to establish database connection, backing off...", "reason", err, "backoff_seconds", backoff.Seconds())
			time.Sleep(backoff)
			backoff = boFormula(backoff)
			try++
		}
	}
}

func (s *AuthService) GetConn(ctx context.Context) *pgx.Conn {
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
	}

	return &conn
}

func (s *AuthService) setDefaultCharacter(u *pbRepo.User) (*pbRepo.Character, error) {
	req := &pbArmory.ListCharactersRequest{Username: u.GetUsername()}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	available, err := s.Clients.Armory.ListCharacters(ctx, req)
	if err != nil {
		return nil, err
	}

	var selected *pbRepo.Character
	if len(available.Characters) == 0 {
		// dont care, create random character
		hero := &pbRepo.Hero{Id: 1}
		req := &pbArmory.CreateCharacterRequest{Hero: hero, Name: u.GetUsername(), UserId: u.GetId()}
		resp, err := s.Clients.Armory.CreateCharacter(ctx, req)
		if err != nil {
			return nil, err
		}

		selected = resp.GetCharacter()
	} else {
		selected = available.Characters[0]
	}

	// update table via selected.CharacterId
	sql, _, err := sq.
		Update("users").
		Set("selected_character_id", selected.GetCharacterId()).
		Where("id = $2").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	_, err = s.DbPool.Exec(ctx, sql, selected.GetCharacterId(), u.GetId())
	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	return selected, nil
}

func (s *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	sql, _, err := sq.
		Select("u.*").
		From("users AS u").
		Where("lower(u.username) LIKE $1").
		Limit(1).
		ToSql()

	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	row := s.DbPool.QueryRow(ctx, sql, req.GetUsername())

	u := &pbRepo.User{}
	err = row.Scan(
		&u.Id,
		&u.Username,
		&u.PasswordHash,
		&u.Email,
		&u.SelectedCharacterId,
	)

	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	var character *pbRepo.Character

	if u.SelectedCharacterId == 0 {
		char, err := s.setDefaultCharacter(u)

		if err != nil {
			s.Logger.Error(err)
			return nil, err
		}

		u.SelectedCharacterId = char.GetCharacterId()
		character = char
	} else {
		req := &pbArmory.GetCharacterRequest{CharacterId: u.SelectedCharacterId}

		resp, err := s.Clients.Armory.GetCharacter(ctx, req)
		if err != nil {
			s.Logger.Error(err)
			return nil, err
		}

		u.SelectedCharacterId = resp.GetCharacter().GetCharacterId()
		character = resp.GetCharacter()
	}

	return &pb.LoginResponse{User: u, Success: true, Character: character}, nil
}

func (s *AuthService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	return nil, nil
}

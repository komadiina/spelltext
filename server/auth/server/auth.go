package server

import (
	"context"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	pb "github.com/komadiina/spelltext/proto/auth"
	pbChar "github.com/komadiina/spelltext/proto/char"
	pbRepo "github.com/komadiina/spelltext/proto/repo"
	"github.com/komadiina/spelltext/server/auth/config"
	"github.com/komadiina/spelltext/utils/singleton/logging"
	"google.golang.org/grpc"
)

type Clients struct {
	Character pbChar.CharacterClient
}

type Connections struct {
	Character *grpc.ClientConn
}

type AuthService struct {
	pb.UnimplementedAuthServer
	DbPool      *pgxpool.Pool
	Config      *config.Config
	Logger      *logging.Logger
	Clients     *Clients
	Connections *Connections
}

func (s *AuthService) setDefaultCharacter(u *pbRepo.User, ctx context.Context) (*pbRepo.Character, error) {
	req := &pbChar.ListCharactersRequest{Username: u.GetUsername()}

	available, err := s.Clients.Character.ListCharacters(ctx, req)
	if err != nil {
		return nil, err
	}

	var selected *pbRepo.Character
	if len(available.Characters) == 0 {
		// dont care, create random character
		hero := &pbRepo.Hero{Id: 1}
		name := strings.ToUpper(u.GetUsername()[0:1]) + u.GetUsername()[1:]
		req := &pbChar.CreateCharacterRequest{Hero: hero, Name: name, UserId: u.GetId()}

		resp, err := s.Clients.Character.CreateCharacter(ctx, req)
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
		char, err := s.setDefaultCharacter(u, ctx)

		if err != nil {
			s.Logger.Error(err)
			return nil, err
		}

		u.SelectedCharacterId = char.GetCharacterId()
		character = char
	} else {
		req := &pbChar.GetCharacterRequest{CharacterId: u.SelectedCharacterId}

		resp, err := s.Clients.Character.GetCharacter(ctx, req)
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

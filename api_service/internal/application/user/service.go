package user

import (
	"api-service/internal/application/ports"
	"api-service/internal/domain"
	userDomain "api-service/internal/domain"
	"context"
	"errors"
	"log/slog"
)

type UserService struct {
	repo         ports.UserRepository
	tokenService ports.TokenService
}

var ErrInvalidCredentials = errors.New("invalid credentials")

func NewUserService(repo ports.UserRepository, tokenService ports.TokenService) *UserService {
	return &UserService{repo: repo, tokenService: tokenService}
}

func (s *UserService) Registration(ctx context.Context, name, email, password string) (string, error) {
	user := domain.NewUser(name, email, password)
	if err := s.repo.Create(ctx, user); err != nil {
		if errors.Is(err, ports.ErrUserAlreadyExists) {
			return "", ports.ErrUserAlreadyExists
		}
		if errors.Is(err, ports.ErrInternal) {
			return "", err
		}
		slog.Error("undefined error", "error", err.Error())
		return "", ports.ErrInternal
	}
	if ctx.Err() != nil {
		return "", context.DeadlineExceeded
	}
	if token, err := s.tokenService.Generate(user.ID()); err != nil {
		slog.Error("token generation error", "error", err.Error())
		return "", ports.ErrInternal
	} else {
		return token, nil
	}
}

func (s *UserService) SignIn(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, ports.ErrUserNotFound) {
			return "", ports.ErrUserNotFound
		} else {
			return "", ports.ErrInternal
		}
	}
	if ctx.Err() != nil {
		return "", context.DeadlineExceeded
	}
	if !userDomain.CheckPasswordHash(password, user.Password()) {
		return "", ErrInvalidCredentials
	}
	tokenString, err := s.tokenService.Generate(user.ID())
	if err != nil {
		slog.Error("token generation error", "error", err.Error())
		return "", ports.ErrInternal
	}
	return tokenString, nil
}

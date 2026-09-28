package service

import (
	"context"

	"github.com/boginskiy/psychologistAI/internal/adapters/dto"
	domain "github.com/boginskiy/psychologistAI/internal/domain/user"
)

type AuthService interface {
	Refresh(context.Context, *dto.TokenReq) (*dto.TokenPair, error)
	Auth(context.Context, *dto.TokenReq) (*dto.InfoUser, error)
	Login(context.Context, *dto.LoginUser) (*dto.TokenPair, error)
}

type UserService interface {
	Verification(ctx context.Context, token string) (*domain.User, error)
	Create(context.Context, *dto.CreateUser) (*domain.User, error)
}

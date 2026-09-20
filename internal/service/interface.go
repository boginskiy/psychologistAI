package service

import (
	"context"

	"github.com/boginskiy/psychologistAI/internal/adapters/dto"
	models "github.com/boginskiy/psychologistAI/internal/models/user"
)

type Validater interface {
	CheckNotEmptyStrField(name, field string) error
}

type Notifier interface {
	Send(email, token string)
}

type AuthService interface {
	Login(context.Context, *dto.LoginUser) (*dto.TokenPair, error)
	Refresh(context.Context, *dto.RefreshTokenRequest) (*dto.TokenPair, error)
}

type UserService interface {
	Create(context.Context, *dto.CreateUser) (*models.User, error)
	Verification(ctx context.Context, token string) (*models.User, error)
}

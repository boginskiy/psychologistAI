package service

import (
	"context"

	"github.com/boginskiy/psychologistAI/internal/adapters/dto"
	models "github.com/boginskiy/psychologistAI/internal/models/users"
)

type Validater interface {
	CheckNotEmptyStrField(name, field string) error
}

type Notifier interface {
	Send(email, token string)
}

type AuthService interface {
	Login(ctx context.Context, loginUser *dto.LoginUser) (*models.Token, error)
	Refresh(ctx context.Context, refreshTokenReq *dto.RefreshTokenRequest) (*models.Token, error)
}

type UserService interface {
	Create(ctx context.Context, user *dto.CreateUser) (*models.User, error)
	Verification(ctx context.Context, token string) (*models.User, error)
}

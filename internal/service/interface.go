package service

import (
	"context"

	"github.com/boginskiy/psychologistAI/internal/adapters/dto"
	"github.com/boginskiy/psychologistAI/internal/models"
)

type Validater interface {
	CheckNotEmptyStrField(name, field string) error
}

type Notifier interface {
	Send(email, token string)
}

type UserService interface {
	Create(ctx context.Context, user *dto.CreateUser) (*models.User, error)
	Verification(ctx context.Context, token string) (*models.User, error)
	Login(ctx context.Context, loginUser *dto.LoginUser) (string, error)
}

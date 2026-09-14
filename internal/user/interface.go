package user

import (
	"context"

	"github.com/boginskiy/psychologistAI/internal/adapters/dto"
	"github.com/boginskiy/psychologistAI/internal/user/models"
)

type UserService interface {
	Create(ctx context.Context, user *dto.CreateUser) (*models.User, error)
	Verification(ctx context.Context, token string) (*models.User, error)
	Login(ctx context.Context, loginUser *dto.LoginUser) (*models.Token, error)
}

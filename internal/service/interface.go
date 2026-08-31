package service

import (
	"context"

	"github.com/boginskiy/psychologistAI/internal/adapters/dto"
)

type Validater interface {
	CheckNotEmptyStrField(name, field string) error
}

type UserService interface {
	CreateUser(ctx context.Context, userReq *dto.CreateUserRequest) (*dto.UserResponse, error)
}

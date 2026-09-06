package service

import (
	"context"

	"github.com/boginskiy/psychologistAI/internal/adapters/dto"
	"github.com/boginskiy/psychologistAI/internal/api/response"
)

type Validater interface {
	CheckNotEmptyStrField(name, field string) error
}

type Notifier interface {
	Send(email, token string)
}

type UserService interface {
	Create(ctx context.Context, userReq *dto.CreateUserRequest) (*response.UserResponse, error)
	Verification(ctx context.Context, token string) (*response.UserResponse, error)
}

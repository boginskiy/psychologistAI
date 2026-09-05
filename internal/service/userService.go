package service

import (
	"context"

	"github.com/boginskiy/psychologistAI/internal/adapters"
	"github.com/boginskiy/psychologistAI/internal/adapters/dto"
	"github.com/boginskiy/psychologistAI/internal/api/response"
	"github.com/boginskiy/psychologistAI/internal/models"
	"github.com/boginskiy/psychologistAI/internal/repository"
)

type UserServ struct {
	Validater Validater
	UserRepo  repository.UserRepo
}

func NewUserServ(ctx context.Context, validater Validater, userRepo repository.UserRepo) *UserServ {
	return &UserServ{
		Validater: validater,
		UserRepo:  userRepo,
	}
}

func (s *UserServ) CreateUser(ctx context.Context, userReq *dto.CreateUserRequest) (*response.UserResponse, error) {
	// Валидация Email
	err := s.Validater.CheckNotEmptyStrField("email", userReq.Email)
	if err != nil {
		return nil, err
	}

	// Валидация Password
	err = s.Validater.CheckNotEmptyStrField("password", userReq.Password)
	if err != nil {
		return nil, err
	}

	// Create domain user
	newUser, err := models.NewUser(userReq)
	if err != nil {
		return nil, err
	}

	// Сохранили user в БД
	err = s.UserRepo.SaveItem(newUser)
	if err != nil {
		return nil, err
	}

	// Выдаем токен для последующей аутентификации

	// TODO
	// БД // Контейнер
	// Выдать токен // теория, куки и т.п.
	// Отправить ответ

	return adapters.ToUserResponse(newUser), nil
}

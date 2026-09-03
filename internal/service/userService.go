package service

import (
	"context"
	"fmt"

	"github.com/boginskiy/psychologistAI/internal/adapters"
	"github.com/boginskiy/psychologistAI/internal/adapters/dto"
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

func (s *UserServ) CreateUser(ctx context.Context, userReq *dto.CreateUserRequest) (*dto.UserResponse, error) {
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

	if !s.UserRepo.CheckUnic(newUser) {
		return nil, fmt.Errorf("user's email is not unique, try again")
	}

	// Не нужно проверки, быстрее будет сохранить пользователя
	// Можем сразу сохранять! И если ошибка уникальности по полю, то выдаем соответствующую ошибку.

	// TODO Сохранить в БД // Контейнер
	// Выдать токен // теория, куки и т.п.
	// Отправить ответ

	return adapters.ToUserResponse(newUser), nil
}

package service

import (
	"context"
	"fmt"

	"github.com/boginskiy/psychologistAI/internal/adapters"
	"github.com/boginskiy/psychologistAI/internal/adapters/dto"
	"github.com/boginskiy/psychologistAI/internal/api/response"
	"github.com/boginskiy/psychologistAI/internal/models"
	"github.com/boginskiy/psychologistAI/internal/repository"
)

type UserServ struct {
	Validater Validater
	Notifier  Notifier
	UserRepo  repository.UserRepo
}

func NewUserServ(ctx context.Context, validater Validater, notifier Notifier, userRepo repository.UserRepo) *UserServ {
	return &UserServ{
		Validater: validater,
		Notifier:  notifier,
		UserRepo:  userRepo,
	}
}

func (s *UserServ) Verification(ctx context.Context, token string) (*response.UserResponse, error) {
	user, err := s.UserRepo.GetItem(token)
	if err != nil {
		return nil, err
	}

	// Обновляем данные после успешной верификации
	user.EmailVerified = true
	user.VerificationToken = ""
	user.TokenExpiresAt = nil

	// Обновляем данные пользователя
	s.UserRepo.UpdateItem(user)

	msg := "verification was successful"
	return adapters.ToUserResponseOnlyMess(msg), nil

}

func (s *UserServ) Create(ctx context.Context, userReq *dto.CreateUserRequest) (*response.UserResponse, error) {
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

	// Отправка асинхронно email для верификации пользователя
	s.Notifier.Send(newUser.Email, newUser.VerificationToken)

	msg := fmt.Sprintf(
		"go to '%s' and verify the account for %v minutes",
		newUser.Email, models.TokenLifetime.Minutes())
	return adapters.ToUserResponseOnlyMess(msg), nil
}

// TODO
// После верификации нужно отправить JWT
// Перекинуть пользователя на анкету, стартовую страницу.

// Что делать с неверифицированными пользователями?
// Ситуация, когда запись добавлена, но письмо не пришло, нужна повторная инициация верификации

// Пользователь переходит по ссылке, сервер проверяет токен, меняет статус на active и, опционально,
// сразу выдаёт JWT (автоматический логин).
// Письма с подтверждением можно пересылать (повторно генерировать новый токен, аннулируя старый).

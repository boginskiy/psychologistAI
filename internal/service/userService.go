package service

import (
	"context"
	"fmt"
	"time"

	"github.com/boginskiy/psychologistAI/internal/adapters"
	"github.com/boginskiy/psychologistAI/internal/adapters/dto"
	"github.com/boginskiy/psychologistAI/internal/api/response"
	"github.com/boginskiy/psychologistAI/internal/models"
	"github.com/boginskiy/psychologistAI/internal/repository"
	"github.com/boginskiy/psychologistAI/pkg/hashpass"
)

const AttemptsCnt = 5

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

// TODO. Слабое место для атак методом перебора.
func (s *UserServ) Verification(ctx context.Context, token string) (*response.UserResponse, error) {
	// Take user from DB
	user, err := s.UserRepo.GetItem(hashpass.CreateHashSHA256(token))
	if err != nil {
		// Need wrap
		return nil, fmt.Errorf("link is incorrect, please try again")
	}

	// Проверка, что EmailVerified == true, т.е. верификация случилась
	if user.EmailVerified == true {
		return nil, fmt.Errorf("client has passed verification")
	}

	// Проверка количеств попыток, данные для верификации.
	if user.Attempts == AttemptsCnt {
		return nil, fmt.Errorf("attempts at verification have run out")
	}
	user.Attempts += 1

	// Check time live of token
	if user.TokenExpiresAt.Before(time.Now().UTC()) {
		verificToken, err := user.UpdateVerificationToken()
		if err != nil {
			return nil, err
		}
		s.UserRepo.UpdateItem(user)               // Update user
		s.Notifier.Send(user.Email, verificToken) // Send email to user

		return nil, fmt.Errorf("link has expired, please try again")
	}

	// User update after good verification
	user.EmailVerified = true
	user.VerificationToken = ""
	user.TokenExpiresAt = nil

	timeNow := time.Now().UTC()
	user.UpdatedAt = &timeNow
	user.VerifiedAt = &timeNow

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

	// Create token
	verificToken, err := newUser.UpdateVerificationToken()
	if err != nil {
		return nil, err
	}

	// Save user in DB
	err = s.UserRepo.SaveItem(newUser)
	if err != nil {
		return nil, err
	}

	// Send email to user
	s.Notifier.Send(newUser.Email, verificToken)

	msg := fmt.Sprintf(
		"go to '%s' and verify the account for %v minutes",
		newUser.Email, models.TokenLifetime.Minutes())
	return adapters.ToUserResponseOnlyMess(msg), nil
}

// Сценарий:
// Учетка не верифицирована
// Неверный токен - что делаем ?
// Просрочено время - что делаем?
// Предусмотреть ограничитель
// Удаление зомби-записи

// Сессионные куки ?

// Если верификация не прошла нужно перенести пользака на /login
// далее - у него есть не верифицированная учетка
// повтороне письмо, учетки нет, перекидываем на регистрацию.

// Учетка верифицирована
// что дальше ?

// Запрос на повторную верификацию. Нужно обновить время, токен
// Soft Delete или Hard Delete -> cron, планировщик базы

// Воркер сам генерирует новый токен и отправляет напоминание («Вы забыли подтвердить регистрацию...»

// TODO
// После верификации нужно отправить JWT
// Перекинуть пользователя на анкету, стартовую страницу.

// Только теперь, когда пользователь помечен как verified, вы генерируете ему основные ключи доступа:

// Access Token (короткоживущий): Например, JWT на 15 минут. Содержит ID пользователя и роль.
// Refresh Token (долгоживущий): Случайная строка, которая сохраняется в базу привязанной к пользователю. Отдается клиенту (лучше всего в HttpOnly куки).

// мидлварь
// auth
// JWT

// Что делать с неверифицированными пользователями?
// Ситуация, когда запись добавлена, но письмо не пришло, нужна повторная инициация верификации

// Пользователь переходит по ссылке, сервер проверяет токен, меняет статус на active и, опционально,
// сразу выдаёт JWT (автоматический логин).
// Письма с подтверждением можно пересылать (повторно генерировать новый токен, аннулируя старый).

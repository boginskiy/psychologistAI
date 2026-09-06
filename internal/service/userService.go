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
	"github.com/boginskiy/psychologistAI/pkg/generators"
	"github.com/boginskiy/psychologistAI/pkg/hashpass"
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
	// hashToken, err := hashpass.CreateHashPass(token)
	// if err != nil {
	// 	return nil, err
	// }

	// bcrypt ?

	// TODO
	// Остановка на хешировании токена верификации

	// Take user from DB
	user, err := s.UserRepo.GetItem(hashToken)
	if err != nil {
		return nil, err
	}

	// Проверка, что EmailVerified == true, т.е. верификация случилась
	if user.EmailVerified == true {
		return nil, fmt.Errorf("client has passed verification")
	}

	// Проверка подлинности токена
	if err := hashpass.CheckPassword(user.VerificationToken, token); err != nil {
		return nil, fmt.Errorf("token is invalid, please try again")
	}

	// if user.VerificationToken != token {
	// 	return nil, fmt.Errorf("token is invalid, please try again")
	// }

	// Проверка время жизни токена
	if user.TokenExpiresAt.After(time.Now().UTC()) {
		return nil, fmt.Errorf("token's time is invalid, please try again")
	}

	// Обновляем данные после успешной верификации
	user.EmailVerified = true
	user.VerificationToken = ""
	user.TokenExpiresAt = nil

	timeNow := time.Now().UTC()
	user.UpdatedAt = &timeNow
	user.VerifiedAt = &timeNow

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

	// Generate Verification Token
	verificToken, err := generators.GenerateToken(models.TokenLength)
	if err != nil {
		return nil, err
	}

	userReq.VerificationToken = verificToken

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
// >>

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

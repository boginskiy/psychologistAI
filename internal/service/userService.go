package service

import (
	"context"
	"fmt"
	"time"

	"github.com/boginskiy/psychologistAI/internal/adapters/dto"
	"github.com/boginskiy/psychologistAI/internal/errs/users"
	models "github.com/boginskiy/psychologistAI/internal/models/users"
	"github.com/boginskiy/psychologistAI/internal/repository"

	"github.com/boginskiy/psychologistAI/pkg/hashpass"
	"github.com/boginskiy/psychologistAI/pkg/jwtservice"
)

const AttemptsCnt = 5

type UserServ struct {
	Validater  Validater
	Notifier   Notifier
	UserRepo   repository.UserRepo
	JWTManager jwtservice.JWTManager
}

func NewUserServ(
	ctx context.Context,
	validater Validater,
	notifier Notifier,
	userRepo repository.UserRepo,
	jwtManager jwtservice.JWTManager,
) *UserServ {
	return &UserServ{
		Validater:  validater,
		Notifier:   notifier,
		UserRepo:   userRepo,
		JWTManager: jwtManager,
	}
}

func (s *UserServ) Login(ctx context.Context, loginUser *dto.LoginUser) (*models.Token, error) {
	// Check user
	userDomain, err := s.UserRepo.GetItem2(loginUser.Email)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", users.ErrInvalidCredentials, err)
	}

	// Check password
	err = hashpass.CheckBcryptPassword(userDomain.HashPassword, loginUser.Password)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", users.ErrInvalidCredentials, err)
	}

	// Verification
	if !userDomain.CheckVerification() {
		if userDomain.Attempts >= AttemptsCnt {
			return nil, users.ErrAttemptsVerification
		}
		userDomain.Attempts += 1
		verificToken, err := userDomain.UpdateVerificationToken()
		if err != nil {
			return nil, fmt.Errorf("%w: %w", users.ErrServer, err)
		}
		s.UserRepo.UpdateItem(userDomain)               // Update user
		s.Notifier.Send(userDomain.Email, verificToken) // Send email to user
		return nil, users.ErrVerification
	}

	// JWT. Generation Access Token
	claim := jwtservice.NewDefaultClaims(
		userDomain.ID,
		userDomain.Role,
		userDomain.Name)

	accessToken, err := s.JWTManager.GenerateToken(claim)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", users.ErrServer, err)
	}

	// Generation Refresh Token
	refreshToken, err := userDomain.UpdateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", users.ErrServer, err)
	}

	// Update user
	s.UserRepo.UpdateItem(userDomain)

	// Create token
	token := models.NewToken(accessToken, refreshToken, claim.GetExpiresIn())

	// Разобрать куки и передать туда refreshToken
	// Иметь ввиду, что при обновлении refreshToken, нужно будет как то находить юзера и
	// вынимать соль для генерации hash и сравненения и уже после генерировать JWT

	return token, nil
}

// TODO. Слабое место для атак методом перебора.
func (s *UserServ) Verification(ctx context.Context, token string) (*models.User, error) {
	// Take user from DB
	userDomain, err := s.UserRepo.GetItem(hashpass.CreateBytesHashSHA256(token))
	if err != nil {
		return nil, fmt.Errorf("%w: %w", users.ErrLinkVerification, err)
	}

	// Проверка, что EmailVerified == true, т.е. верификация случилась
	if userDomain.EmailVerified == true {
		return nil, users.ErrRepeatVerification
	}

	// Проверка количеств попыток, данные для верификации.
	if userDomain.Attempts >= AttemptsCnt {
		return nil, users.ErrAttemptsVerification
	}
	userDomain.Attempts += 1

	// Check time live of token
	if userDomain.TokenExpiresAt.Before(time.Now().UTC()) {
		verificToken, err := userDomain.UpdateVerificationToken()
		if err != nil {
			return nil, fmt.Errorf("%w: %w", users.ErrServer, err)
		}
		s.UserRepo.UpdateItem(userDomain)               // Update user
		s.Notifier.Send(userDomain.Email, verificToken) // Send email to user

		return nil, users.ErrVerification
	}

	// User update after good verification
	userDomain.EmailVerified = true
	userDomain.HashVerifToken = []byte{}
	userDomain.TokenExpiresAt = nil

	timeNow := time.Now().UTC()
	userDomain.UpdatedAt = &timeNow
	userDomain.VerifiedAt = &timeNow

	s.UserRepo.UpdateItem(userDomain)

	return userDomain, nil
}

func (s *UserServ) Create(ctx context.Context, createUser *dto.CreateUser) (*models.User, error) {
	// Валидация Email
	err := s.Validater.CheckNotEmptyStrField("email", createUser.Email)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", users.ErrInvalidCredentials, err)
	}

	// Валидация Password
	err = s.Validater.CheckNotEmptyStrField("password", createUser.Password)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", users.ErrInvalidCredentials, err)
	}

	// Create domain user
	userDomain, err := models.NewUser(createUser)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", users.ErrServer, err)
	}

	// Create token
	verificToken, err := userDomain.UpdateVerificationToken()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", users.ErrServer, err)
	}

	// Save user in DB
	err = s.UserRepo.SaveItem(userDomain)
	if err != nil {
		// TODО, пока отправляем ошибку сервера, но в целом у пользователя может быть не уникальный email
		// и тогда ему надо что то передать.
		return nil, fmt.Errorf("%w: %w", users.ErrServer, err)
	}

	// Send email to ErrServeruser
	s.Notifier.Send(userDomain.Email, verificToken)

	return userDomain, nil
}

// Сценарий:
// Предусмотреть кнопку "выйти"

// Удаление зомби-записи

// Сессионные куки ?

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

// Пользователь переходит по ссылке, сервер проверяет токен, меняет статус на active и, опционально,
// сразу выдаёт JWT (автоматический логин).

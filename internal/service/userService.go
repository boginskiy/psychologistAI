package service

import (
	"context"
	"fmt"
	"time"

	"github.com/boginskiy/psychologistAI/internal/adapters/dto"
	"github.com/boginskiy/psychologistAI/internal/errs/server"
	"github.com/boginskiy/psychologistAI/internal/errs/users"
	models "github.com/boginskiy/psychologistAI/internal/models/user"
	"github.com/boginskiy/psychologistAI/internal/repository"

	"github.com/boginskiy/psychologistAI/pkg/hashpass"
	"github.com/boginskiy/psychologistAI/pkg/jwtservice"
)

// const AttemptsCnt = 5

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

// TODO. Слабое место для атак методом перебора.
func (s *UserServ) Verification(ctx context.Context, token string) (*models.User, error) {
	// Take user from DB
	userDomain, err := s.UserRepo.ReadByToken(hashpass.CreateBytesHashSHA256(token))
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
	if userDomain.ExpiresAtVerifToken.Before(time.Now().UTC()) {
		verificToken, err := userDomain.UpdateVerificationToken()
		if err != nil {
			return nil, fmt.Errorf("%w: %w", server.ErrServer, err)
		}
		s.UserRepo.UpdateItem(userDomain)               // Update user
		s.Notifier.Send(userDomain.Email, verificToken) // Send email to user

		return nil, users.ErrVerification
	}

	// User update after good verification
	userDomain.EmailVerified = true
	userDomain.HashVerifToken = []byte{}
	userDomain.ExpiresAtVerifToken = nil

	timeNow := time.Now().UTC()
	userDomain.UpdatedAt = &timeNow
	userDomain.VerifiedAtVerifToken = &timeNow

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
		return nil, fmt.Errorf("%w: %w", server.ErrServer, err)
	}

	// Create token
	verificToken, err := userDomain.UpdateVerificationToken()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", server.ErrServer, err)
	}

	// Save user in DB
	err = s.UserRepo.Create(userDomain)
	if err != nil {
		// TODО, пока отправляем ошибку сервера, но в целом у пользователя может быть не уникальный email
		// и тогда ему надо что то передать.
		return nil, fmt.Errorf("%w: %w", server.ErrServer, err)
	}

	// Send email to ErrServeruser
	s.Notifier.Send(userDomain.Email, verificToken)

	return userDomain, nil
}

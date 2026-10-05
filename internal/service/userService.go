package service

import (
	"context"
	"fmt"
	"time"

	"github.com/boginskiy/psychologistAI/internal/adapters/dto"
	"github.com/boginskiy/psychologistAI/internal/errs"

	domain "github.com/boginskiy/psychologistAI/internal/domain/user"
	"github.com/boginskiy/psychologistAI/internal/repository"

	"github.com/boginskiy/psychologistAI/pkg/emailservice"
	"github.com/boginskiy/psychologistAI/pkg/hashpass"
	"github.com/boginskiy/psychologistAI/pkg/jwtservice"
	"github.com/boginskiy/psychologistAI/pkg/validservice"
)

// const AttemptsCnt = 5

type UserServ struct {
	Validater  validservice.Validater
	Postman    emailservice.Postman
	UserRepo   repository.UserRepo
	JWTManager jwtservice.JWTManager
}

func NewUserServ(
	ctx context.Context,
	validater validservice.Validater,
	postman emailservice.Postman,
	userRepo repository.UserRepo,
	jwtManager jwtservice.JWTManager,
) *UserServ {
	return &UserServ{
		Validater:  validater,
		Postman:    postman,
		UserRepo:   userRepo,
		JWTManager: jwtManager,
	}
}

// TODO. Слабое место для атак методом перебора.
func (s *UserServ) Verification(ctx context.Context, token string) (*domain.User, error) {
	// Take user from DB
	userDomain, err := s.UserRepo.ReadByToken(hashpass.CreateBytesHashSHA256(token))
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errs.ErrLinkVerification, err)
	}

	// Проверка, что EmailVerified == true, т.е. верификация случилась
	if userDomain.EmailVerified == true {
		return nil, errs.ErrRepeatVerification
	}

	// Проверка количеств попыток, данные для верификации.
	if userDomain.Attempts >= Attempts {
		return nil, errs.ErrAttemptsVerification
	}
	userDomain.Attempts += 1

	// Check time live of token
	if userDomain.ExpiresAtVerifToken.Before(time.Now().UTC()) {
		verificToken, err := userDomain.UpdateVerificationToken()
		if err != nil {
			return nil, fmt.Errorf("%w: %w", errs.ErrServer, err)
		}

		// Update user
		s.UserRepo.UpdateItem(userDomain)

		// Create message and Send email to user
		verifLetter := emailservice.NewVerifLetter(userDomain.Email, verificToken, NameProject)
		s.Postman.Send(verifLetter, Retry) // Send email to user

		return nil, errs.ErrVerification
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

func (s *UserServ) Create(ctx context.Context, createUser *dto.CreateUser) (*domain.User, error) {
	// Валидация Email
	err := s.Validater.CheckNotEmptyStrField("email", createUser.Email)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errs.ErrInvalidCredentials, err)
	}

	// Валидация Password
	err = s.Validater.CheckNotEmptyStrField("password", createUser.Password)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errs.ErrInvalidCredentials, err)
	}

	// Create domain user
	userDomain, err := domain.NewUser(createUser)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errs.ErrServer, err)
	}

	// Create token
	verificToken, err := userDomain.UpdateVerificationToken()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errs.ErrServer, err)
	}

	// Save user in DB
	err = s.UserRepo.Create(userDomain)
	if err != nil {
		// TODО, пока отправляем ошибку сервера, но в целом у пользователя может быть не уникальный email
		// и тогда ему надо что то передать.
		return nil, fmt.Errorf("%w: %w", errs.ErrRegistr, err)
	}

	// Path for verification user
	verificationPath := createUser.VerificationLink + verificToken

	// Create message and Send email to user
	verifLetter := emailservice.NewVerifLetter(userDomain.Email, verificationPath, NameProject)
	s.Postman.Send(verifLetter, Retry)

	return userDomain, nil
}

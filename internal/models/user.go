package models

import (
	"time"

	"github.com/boginskiy/psychologistAI/internal/adapters/dto"
	"github.com/boginskiy/psychologistAI/pkg/hashpass"
	"github.com/google/uuid"
)

const TokenLength = 32
const TokenLifetime = 15 * time.Minute

type User struct {
	// Основные поля
	ID       uuid.UUID `json:"id" db:"id"`
	Email    string    `json:"email" db:"email" validate:"required,email"`
	Password string    `json:"-" db:"password_hash"` // Храним хеш, не выводим

	// Личная информация
	FirstName string `json:"first_name,omitempty" db:"first_name"`
	LastName  string `json:"last_name,omitempty" db:"last_name"`
	Phone     string `json:"phone,omitempty" db:"phone"`

	// Статусы
	IsActive bool   `json:"is_active" db:"is_active"`
	IsAdmin  bool   `json:"is_admin" db:"is_admin"`
	Role     string `json:"role" db:"role"` // user, admin, moderator

	// Временные метки
	CreatedAt  *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at" db:"updated_at"`
	VerifiedAt *time.Time `json:"verified_at" db:"verified_at"`

	LastActivityAt *time.Time `json:"last_activity_at,omitempty" db:"last_activity_at"`
	DeletedAt      *time.Time `json:"-" db:"deleted_at"` // soft delete

	// Токен верификации пользователя
	TokenExpiresAt    *time.Time `json:"token_expires_at" db:"token_expires_at"`
	EmailVerified     bool       `json:"email_verified" db:"email_verified"`
	VerificationToken string     `json:"verification_token" db:"verification_token"`
}

func NewUser(userReq *dto.CreateUserRequest) (*User, error) {
	hashPassword, err := hashpass.CreateHashPass(userReq.Password)
	if err != nil {
		return nil, err
	}

	hashVerificationToken, err := hashpass.CreateHashPass(userReq.VerificationToken)
	if err != nil {
		return nil, err
	}

	// // Токен для верификации пользователя
	// verificToken, err := generators.GenerateToken(TokenLength)
	// if err != nil {
	// 	return nil, err
	// }

	// Время жизни токена до N-time
	tokenExpiresAt := time.Now().UTC().Add(TokenLifetime)

	// Current time
	timeNow := time.Now().UTC()

	return &User{
		ID:                uuid.Must(uuid.NewV7()),
		Email:             userReq.Email,
		Password:          hashPassword,
		FirstName:         userReq.FirstName,
		LastName:          userReq.LastName,
		Phone:             userReq.Phone,
		Role:              "user",
		CreatedAt:         &timeNow,
		UpdatedAt:         &timeNow,
		VerifiedAt:        nil,
		EmailVerified:     false,
		VerificationToken: hashVerificationToken,
		TokenExpiresAt:    &tokenExpiresAt,
	}, nil
}

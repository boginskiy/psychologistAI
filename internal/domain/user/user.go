package user

import (
	"crypto/subtle"
	"time"

	"github.com/boginskiy/psychologistAI/cmd/config"
	"github.com/boginskiy/psychologistAI/internal/adapters/dto"
	"github.com/boginskiy/psychologistAI/pkg/generators"
	"github.com/boginskiy/psychologistAI/pkg/hashpass"
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Email        string    `json:"email" db:"email" validate:"required,email"`
	HashPassword string    `json:"-" db:"hash_password"` // Храним хеш, не выводим
	Name         string    `json:"name,omitempty" db:"name"`
	Phone        string    `json:"phone,omitempty" db:"phone"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	IsAdmin      bool      `json:"is_admin" db:"is_admin"`
	Role         []string  `json:"role" db:"role"` // user, admin, moderator

	// Временные метки
	CreatedAt      *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at" db:"updated_at"`
	LastActivityAt *time.Time `json:"last_activity_at,omitempty" db:"last_activity_at"`
	DeletedAt      *time.Time `json:"-" db:"deleted_at"` // soft delete

	// Верификации пользователя.Need Table
	HashVerifToken       []byte     `json:"hash_verif_token" db:"hash_verif_token"`
	ExpiresAtVerifToken  *time.Time `json:"expires_at_verif_token" db:"expires_at_verif_token"`
	VerifiedAtVerifToken *time.Time `json:"verified_at_verif_token" db:"verified_at_verif_token"`
	EmailVerified        bool       `json:"email_verified" db:"email_verified"`
	Attempts             int        `json:"attempts" db:"attempts"`
}

func NewUser(createUser *dto.CreateUser) (*User, error) {
	hashPassword, err := hashpass.CreateBcryptHashPassword(createUser.Password)
	if err != nil {
		return nil, err
	}
	// Current time
	timeNow := time.Now().UTC()

	return &User{
		ID:                   generators.CreateUUIDv7(),
		Email:                createUser.Email,
		HashPassword:         hashPassword,
		Name:                 createUser.Name,
		Phone:                createUser.Phone,
		Role:                 []string{"user"},
		CreatedAt:            &timeNow,
		UpdatedAt:            &timeNow,
		VerifiedAtVerifToken: nil,
		EmailVerified:        false,
	}, nil
}

func (u *User) CheckVerification() bool {
	return u.EmailVerified && u.VerifiedAtVerifToken != nil
}

func (u *User) UpdateVerificationToken() (string, error) {
	// Generate New Verification Token
	verificToken, err := generators.GenerateTokenBase64(config.LENGTH_VARIFICATION_TOKEN)
	if err != nil {
		return "", err
	}
	u.HashVerifToken = hashpass.CreateBytesHashSHA256(verificToken)

	timeMinute := time.Duration(config.LIVE_TIME_VARIFICATION_TOKEN) * time.Minute
	tokenExpiresAt := time.Now().UTC().Add(timeMinute)

	u.ExpiresAtVerifToken = &tokenExpiresAt
	return verificToken, nil
}

func (u *User) CompareHash(hashToken []byte) bool {
	return subtle.ConstantTimeCompare(u.HashVerifToken, hashToken) == 1
}

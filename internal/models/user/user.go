package models

import (
	"crypto/subtle"
	"fmt"
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

	// Верификации пользователя
	HashVerifToken       []byte     `json:"hash_verif_token" db:"hash_verif_token"`
	ExpiresAtVerifToken  *time.Time `json:"expires_at_verif_token" db:"expires_at_verif_token"`
	VerifiedAtVerifToken *time.Time `json:"verified_at_verif_token" db:"verified_at_verif_token"`
	EmailVerified        bool       `json:"email_verified" db:"email_verified"`
	Attempts             int        `json:"attempts" db:"attempts"`

	// Аутентификация пользователя
	HashRefreshToken      []byte     `json:"hash_refresh_token" db:"hash_refresh_token"`
	Salt                  []byte     `json:"salt" db:"salt"`
	ExpiresAtRefreshToken *time.Time `json:"expires_at_refresh_token" db:"expires_at_refresh_token"`
	DeviceInfo            DeviceInfo `json:"device_info" db:"device_info"`
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
	tokenExpiresAt := time.Now().UTC().Add(config.LIVE_TIME_VARIFICATION_TOKEN)
	u.ExpiresAtVerifToken = &tokenExpiresAt
	return verificToken, nil
}

func (u *User) CompareHash(hashToken []byte) bool {
	return subtle.ConstantTimeCompare(u.HashVerifToken, hashToken) == 1
}

func (u *User) UpdateRefreshToken(ip, userAgent string) (string, error) {
	// err := u.updateSalt()
	// if err != nil {
	// 	return "", err
	// }

	// TODO? Зашить в токен данные! Например ID пользователя

	token := generators.CreateUUIDv7ToString()
	tokenFull := u.assembleFullToken(token, ip, userAgent)

	u.HashRefreshToken = hashpass.CreateBytesHashSHA256WithSalt(u.Salt, tokenFull)
	u.updateRefreshExpiresAt()

	return token, nil
}

func (u *User) assembleFullToken(token, ip, userAgent string) string {
	return fmt.Sprintf("%s:%s:%s", token, ip, userAgent)
}

func (u *User) updateRefreshExpiresAt() {
	refreshTokenExpiresAt := time.Now().UTC().Add(config.LIVE_TIME_REFRESH_TOKEN)
	u.ExpiresAtRefreshToken = &refreshTokenExpiresAt
}

// func (u *User) updateSalt() error {
// 	salt, err := generators.GenerateRandomBytes(config.LENGTH_SALT)
// 	if err != nil {
// 		return err
// 	}
// 	u.Salt = salt
// 	return nil
// }

// // VerifyRefreshToken проверяет входящий токен.
// func VerifyRefreshToken(incomingToken string, storedSalt []byte, storedHash []byte) bool {
// 	hasher := sha256.New()
// 	hasher.Write(storedSalt)
// 	hasher.Write([]byte(incomingToken))
// 	computedSum := hasher.Sum(nil)

// 	// Используем ConstantTimeCompare, чтобы избежать атак по времени (timing attacks)
// 	// Никогда не используйте просто == для сравнения хешей!
// 	return subtle.ConstantTimeCompare(computedSum, storedHash) == 1
// }

package models

import (
	"crypto/subtle"
	"time"

	"github.com/boginskiy/psychologistAI/internal/adapters/dto"
	"github.com/boginskiy/psychologistAI/pkg/generators"
	"github.com/boginskiy/psychologistAI/pkg/hashpass"
	"github.com/google/uuid"
)

const TokenLength = 32
const VerifTokenLifetime = 15 * time.Minute     // 15 минут
const RefreshTokenLifetime = 24 * 7 * time.Hour // 7 дней
const SaltLength = 16

type User struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Email        string    `json:"email" db:"email" validate:"required,email"`
	HashPassword string    `json:"-" db:"hash_password"` // Храним хеш, не выводим
	Name         string    `json:"name,omitempty" db:"name"`
	Phone        string    `json:"phone,omitempty" db:"phone"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	IsAdmin      bool      `json:"is_admin" db:"is_admin"`
	Role         string    `json:"role" db:"role"` // user, admin, moderator

	// Временные метки
	CreatedAt      *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at" db:"updated_at"`
	LastActivityAt *time.Time `json:"last_activity_at,omitempty" db:"last_activity_at"`
	DeletedAt      *time.Time `json:"-" db:"deleted_at"` // soft delete

	// Верификации пользователя
	HashVerifToken []byte     `json:"hash_verif_token" db:"hash_verif_token"`
	TokenExpiresAt *time.Time `json:"token_expires_at" db:"token_expires_at"`
	VerifiedAt     *time.Time `json:"verified_at" db:"verified_at"`
	EmailVerified  bool       `json:"email_verified" db:"email_verified"`
	Attempts       int        `json:"attempts" db:"attempts"`

	// Аутентификация пользователя
	HashRefreshToken      []byte     `json:"hash_refresh_token" db:"hash_refresh_token"`
	Salt                  []byte     `json:"salt" db:"salt"`
	RefreshTokenExpiresAt *time.Time `json:"refresh_token_expires_at" db:"refresh_token_expires_at"`
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
		ID:            generators.CreateUUIDv7(),
		Email:         createUser.Email,
		HashPassword:  hashPassword,
		Name:          createUser.Name,
		Phone:         createUser.Phone,
		Role:          "user",
		CreatedAt:     &timeNow,
		UpdatedAt:     &timeNow,
		VerifiedAt:    nil,
		EmailVerified: false,
	}, nil
}

func (u *User) CheckVerification() bool {
	return u.EmailVerified && u.VerifiedAt != nil
}

func (u *User) UpdateVerificationToken() (string, error) {
	// Generate New Verification Token
	verificToken, err := generators.GenerateTokenBase64(TokenLength)
	if err != nil {
		return "", err
	}
	u.HashVerifToken = hashpass.CreateBytesHashSHA256(verificToken)
	tokenExpiresAt := time.Now().UTC().Add(VerifTokenLifetime)
	u.TokenExpiresAt = &tokenExpiresAt
	return verificToken, nil
}

func (u *User) CompareHash(hashToken []byte) bool {
	return subtle.ConstantTimeCompare(u.HashVerifToken, hashToken) == 1
}

func (u *User) UpdateRefreshToken() (string, error) {
	err := u.updateSalt()
	if err != nil {
		return "", err
	}

	u.updateRefreshExpiresAt()

	token := generators.CreateUUIDv7ToString()
	u.HashRefreshToken = hashpass.CreateBytesHashSHA256WithSalt(u.Salt, token)

	return token, nil
}

func (u *User) updateRefreshExpiresAt() {
	refreshTokenExpiresAt := time.Now().UTC().Add(RefreshTokenLifetime)
	u.RefreshTokenExpiresAt = &refreshTokenExpiresAt
}

func (u *User) updateSalt() error {
	salt, err := generators.GenerateRandomBytes(SaltLength)
	if err != nil {
		return err
	}
	u.Salt = salt
	return nil
}

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

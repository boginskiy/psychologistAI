package jwtservice

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

const TIME_LIVE_TOKEN = 15 * time.Minute
const HOST_SITE = "psychologistAI.com"
const JWT_SECRET_KEY = "cjlsjdc3r483ucdhcyeruf9ehrc"

type DefaultClaims struct {
	ID   uuid.UUID `json:"id"`
	Role string    `json:"role"`
	Name string    `json:"name"`
	jwt.RegisteredClaims
	expiresIn int
}

func NewDefaultClaims(id uuid.UUID, role, name string) Claims {
	issuedAt := time.Now()
	expiresAt := issuedAt.Add(TIME_LIVE_TOKEN)

	return &DefaultClaims{
		ID:   id,
		Role: role,
		Name: name,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(issuedAt),  // Время выдачи
			ExpiresAt: jwt.NewNumericDate(expiresAt), // Срок действия: токен станет невалидным через 15 минут
			Issuer:    HOST_SITE,                     // Кто выдал токен (рекомендуется заполнять)
		},
		expiresIn: int(expiresAt.Sub(issuedAt).Seconds()),
	}
}

func (c *DefaultClaims) GetExpiresIn() int { return c.expiresIn }
func (c *DefaultClaims) GetID() uuid.UUID  { return c.ID }
func (c *DefaultClaims) GetRole() string   { return c.Role }
func (c *DefaultClaims) GetName() string   { return c.Name }

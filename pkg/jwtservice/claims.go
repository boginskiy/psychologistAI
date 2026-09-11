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
}

func NewDefaultClaims(id uuid.UUID, role, name string) Claims {
	return &DefaultClaims{
		ID:   id,
		Role: role,
		Name: name,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),                      // Время выдачи
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TIME_LIVE_TOKEN)), // Срок действия: токен станет невалидным через 15 минут
			Issuer:    HOST_SITE,                                           // Кто выдал токен (рекомендуется заполнять)
		},
	}
}

func (c *DefaultClaims) GetID() uuid.UUID { return c.ID }
func (c *DefaultClaims) GetRole() string  { return c.Role }
func (c *DefaultClaims) GetName() string  { return c.Name }

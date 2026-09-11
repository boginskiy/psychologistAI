package jwtservice

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type Claims interface {
	jwt.Claims
	GetID() uuid.UUID
	GetRole() string
	GetName() string
}

type JWTManager interface {
	GenerateToken(claims Claims) (string, error)
	CheckAndParseToken(token string) (Claims, error)
}

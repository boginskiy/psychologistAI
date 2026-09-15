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
	GetExpiresIn() int
}

type JWTManager interface {
	GenerateToken(claims Claims) (string, error)
	CheckAndParseToken(tokenStr string, claims Claims) (Claims, error)
}

package jwtservice

import "github.com/golang-jwt/jwt/v4"

type ClaimConfig interface {
	GetTimeLiveToken() int
	GetHostName() string
}

type JWTManager interface {
	GenerateToken(secretKey string, claim jwt.Claims) (string, error)
	CheckAndParseToken(secretKey string, token string, claim jwt.Claims) (jwt.Claims, error)
}

package jwtservice

import "github.com/golang-jwt/jwt/v4"

type TokenConfig interface {
	GetTimeLiveToken() int
	GetSecretKeyForToken() string
}

type HostConfig interface {
	GetHostName() string
}

type JWTConfig interface {
	TokenConfig
	HostConfig
}

type JWTManager interface {
	GenerateToken(JWTConfig, jwt.Claims) (string, error)
	CheckAndParseToken(JWTConfig, string, jwt.Claims) (jwt.Claims, error)
}

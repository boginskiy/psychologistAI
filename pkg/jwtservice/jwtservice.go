package jwtservice

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

type JWTService struct {
}

func NewJWTService() *JWTService {
	return &JWTService{}
}

func (s *JWTService) GenerateToken(config JWTConfig, claims jwt.Claims) (string, error) {
	// Создаем новый токен с алгоритмом HS256 (HMAC с SHA256)
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Подписываем токен нашим секретным ключом
	key := []byte(config.GetSecretKeyForToken())

	signedToken, err := newToken.SignedString(key)
	if err != nil {
		fmt.Printf("%v: %v", ErrTokenSigning, err) // + logger err
		return "", fmt.Errorf("%w: %w", ErrTokenSigning, err)
	}
	return signedToken, nil
}

func (s *JWTService) CheckAndParseToken(config JWTConfig, tokenString string, claims jwt.Claims) (jwt.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return config.GetSecretKeyForToken(), nil
	})

	if err != nil {
		var vErr *jwt.ValidationError

		if errors.As(err, &vErr) {
			// Проверка истечения срока действия
			if errors.Is(vErr.Inner, jwt.ErrTokenExpired) {
				fmt.Printf("%v: %v", ErrTokenLive, err) // + logger err
				return claims, fmt.Errorf("%w: %w", ErrTokenLive, err)
			}
		}

		fmt.Printf("%v: %v", ErrParceToken, err) // + logger err
		return nil, fmt.Errorf("%w: %w", ErrParceToken, err)
	}

	if !token.Valid {
		fmt.Printf("%v", ErrTokenValid) // + logger err
		return nil, ErrTokenValid
	}

	return claims, nil
}

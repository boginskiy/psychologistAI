package jwtservice

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

type JWTService struct {
}

func NewJWTService() *JWTService {
	return nil
}

func (s *JWTService) GenerateToken(claims Claims) (string, error) {
	// Создаем новый токен с алгоритмом HS256 (HMAC с SHA256)
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Подписываем токен нашим секретным ключом
	signedToken, err := newToken.SignedString(JWT_SECRET_KEY)
	if err != nil {
		// logger err
		return "", fmt.Errorf("%w: %w", ErrTokenSigning, err)
	}
	return signedToken, nil
}

func (s *JWTService) CheckAndParseToken(tokenStr string, claims Claims) (Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return JWT_SECRET_KEY, nil
	})

	if err != nil {
		var vErr *jwt.ValidationError
		// logger err
		if errors.As(err, &vErr) {
			// Проверка истечения срока действия
			if errors.Is(vErr.Inner, jwt.ErrTokenExpired) {
				return claims, fmt.Errorf("%w: %w", ErrTokenLive, err)
			}
		}
		return nil, fmt.Errorf("%w: %w", ErrParceToken, err)
	}

	if !token.Valid {
		// logger err
		return nil, ErrTokenValid
	}

	return claims, nil
}

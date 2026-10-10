package jwtservice

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

const TokenTypeAccess = "access"
const TokenTypeRefresh = "refresh"

type AccessTokenClaim struct {
	UserID               uuid.UUID `json:"uid"`        // Краткий ключ экономит байты в заголовке HTTP
	UserRole             string    `json:"role"`       // Права доступа (RBAC)
	TokenType            string    `json:"type"`       // "access"
	SessionID            string    `json:"session_id"` // jti из refresh token, чтобы связать их
	jwt.RegisteredClaims           // Стандартные поля JWT
}

func NewAccessTokenClaim(config ClaimConfig, refreshClaim *RefreshTokenClaim) *AccessTokenClaim {
	timeMinute := time.Duration(config.GetTimeLiveToken()) * time.Minute
	now := time.Now()

	return &AccessTokenClaim{
		UserID:    refreshClaim.UserID,
		UserRole:  refreshClaim.UserRole,
		TokenType: TokenTypeAccess,
		SessionID: refreshClaim.ID,

		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    config.GetHostName(),
			Subject:   refreshClaim.UserID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(timeMinute)),
			// Audience: []string{"web-client", "mobile-app"},
		},
	}
}

// RefreshTokenClaim Default
type RefreshTokenClaim struct {
	UserID               uuid.UUID `json:"uid"`  // Subject — идентификатор владельца (лучше называть UserID для ясности)
	UserName             string    `json:"name"` // Имя пользователя (для отображения без доп. запросов к БД)
	UserRole             string    `json:"role"` // Роли (Scope/Permissions). Слайс позволяет иметь RBAC.
	TokenType            string    `json:"type"` // "refresh" — чтобы отличать от access token
	jwt.RegisteredClaims           // Вложенная анонимная структура со стандартными полями
}

func NewRefreshTokenClaim(config ClaimConfig, tokenUser *TokenUser) *RefreshTokenClaim {
	timeMinute := time.Duration(config.GetTimeLiveToken()) * time.Minute
	now := time.Now()

	return &RefreshTokenClaim{
		UserID:    tokenUser.ID,
		UserName:  tokenUser.Name,
		UserRole:  tokenUser.Role,
		TokenType: TokenTypeRefresh,

		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.Must(uuid.NewV7()).String(),        // Уникальный номер этого конкретного токена (jti)
			Issuer:    config.GetHostName(),                    // Кто выдал токен (рекомендуется заполнять)
			Subject:   tokenUser.ID.String(),                   // Идентификатор пользователя (sub)
			IssuedAt:  jwt.NewNumericDate(now),                 // Время выдачи
			NotBefore: jwt.NewNumericDate(now),                 // Время начала действия токена
			ExpiresAt: jwt.NewNumericDate(now.Add(timeMinute)), // Срок действия: токен станет невалидным через N минут
			// Audience: []string{"web-client"},
		},
	}
}

package dto

import "time"

// RefreshTokenRequest - DTO для refresh token
type RefreshTokenRequest struct {
	Token     string
	IP        string
	UserAgent string
}

// TokenPair - DTO
type TokenPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	SessionID    string    `json:"session_id"`  // Или "jti" / "sid", если sessionID совпадает с jti refresh-токена
	SessionExp   time.Time `json:"session_exp"` // Время жизни refresh token и соответственно session
}

// Token - DTO
type Token struct {
	Access  string
	Refresh string
}

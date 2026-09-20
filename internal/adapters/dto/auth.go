package dto

// RefreshTokenRequest - DTO для refresh token
type RefreshTokenRequest struct {
	Token     string
	IP        string
	UserAgent string
}

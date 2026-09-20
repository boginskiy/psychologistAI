package response

import models "github.com/boginskiy/psychologistAI/internal/models/user"

const TokenType = "bearer"

type TokenResponse struct {
	AccessToken string `json:"access_token,omitempty"`
	ExpiresIn   int    `json:"expires_in,omitempty"`
	TokenType   string `json:"token_type,omitempty"`
	InfoResponse
}

func (r *TokenResponse) AttrsUpdate(token *models.Token, status int) {
	r.AccessToken = token.Access
	r.ExpiresIn = token.ExpiresIn
	r.TokenType = TokenType
	r.Status = status
}

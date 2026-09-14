package response

const TokenType = "bearer"

type TokenRes struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TType       string `json:"token_type"`
	InfoRes
}

// Разобраться с респонсами! Что как лучше сделать!

// "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...", "expires_in": 900, "token_type": "bearer"

func NewTokenRes(token, tp string, expiresIn int) *TokenRes {
	return &TokenRes{
		AccessToken: token,
		ExpiresIn:   expiresIn,
		TType:       TokenType}
}

func (r *TokenRes) GetStatus() int {
	return r.Status
}

func (r *TokenRes) UpdateErr(err error, status int) {
	r.Message = err.Error()
	r.Status = status
}

func (r *TokenRes) UpdateInfo(info string, status int) {
	r.Message = info
	r.Status = status
}

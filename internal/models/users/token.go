package models

type Token struct {
	Access    string
	ExpiresIn int
	Refresh   string
}

func NewToken(access, refresh string, expiresIn int) *Token {
	return &Token{
		Access:    access,
		Refresh:   refresh,
		ExpiresIn: expiresIn,
	}
}

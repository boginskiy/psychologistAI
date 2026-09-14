package models

type Token struct {
	Access  string
	Refresh string
}

func NewToken(access, refresh string) *Token {
	return &Token{
		Access:  access,
		Refresh: refresh,
	}
}

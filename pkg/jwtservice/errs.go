package jwtservice

import "errors"

var (
	ErrTokenLive    error = errors.New("token has expired")
	ErrTokenValid   error = errors.New("token is invalid")
	ErrTokenSigning error = errors.New("token signing failed")
	ErrParceToken   error = errors.New("token parsing error")
)

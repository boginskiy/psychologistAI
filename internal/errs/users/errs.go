package users

import "errors"

var (
	// Verification
	ErrVerification         = errors.New("account verification has not been completed, check your email")
	ErrAttemptsVerification = errors.New("verification limit has been exceeded")
	ErrLinkVerification     = errors.New("link is incorrect, please try again")
	ErrRepeatVerification   = errors.New("client has passed verification")

	// Credentials
	ErrInvalidCredentials = errors.New("invalid email or password")

	// Sever:
	ErrServer = errors.New("server error")
)

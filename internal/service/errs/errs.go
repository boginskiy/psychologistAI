package errs

import "errors"

var (
	ErrVerification       = errors.New("account verification has not been completed")
	ErrAttemptsVerific    = errors.New("verification limit has been exceeded")
	ErrUpdateVerificToken = errors.New("error updating the verification token")
)

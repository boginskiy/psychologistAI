package errs

import "errors"

var (
	// Verification
	ErrVerification         = errors.New("account verification has not been completed, check your email")
	ErrAttemptsVerification = errors.New("verification limit has been exceeded")
	ErrLinkVerification     = errors.New("link is incorrect, please try again")
	ErrRepeatVerification   = errors.New("client has passed verification")

	// Authentification
	ErrSession         = errors.New("session not found")
	ErrCompareToken    = errors.New("tokens don't match")
	ErrTimeLiveSession = errors.New("session is over")
	ErrAuth            = errors.New("authentication error, need to log in")
	ErrUsingToken      = errors.New("reuse of the token")
	ErrLegitimacyUser  = errors.New("user legitimacy check failed")

	// Session
	ErrHistoricalSession = errors.New("historical session is missing")

	// Credentials
	ErrInvalidCredentials = errors.New("invalid email or password")

	// Server
	ErrServer = errors.New("server error")
)

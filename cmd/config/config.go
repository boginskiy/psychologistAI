package config

const (
	// Server
	HOST_NAME = "psychologistAI.com"

	// Tokens-JWT
	TIME_LIVE_JWT_ACCESS_TOKEN  = 15 // Minute
	SECRET_KEY_JWT_ACCESS_TOKEN = "cjlsjdc3r483ucdhcyeruf9ehrc"

	TIME_LIVE_JWT_REFRESH_TOKEN  = 10080 // Minute
	SECRET_KEY_JWT_REFRESH_TOKEN = "acsdcfwecwX2343FCSDBJHXSDJ"

	// Cookies
	COOKIE_NAME_ACCESS_TOKEN      = "access_token"
	COOKIE_EXPIRES_ACCESS_TOKEN   = 15  // Minute
	COOKIE_MAX_AGE_ACCESS_TOKEN   = 900 // Second
	COOKIE_PATH_ACCESS_TOKEN      = "/"
	COOKIE_HTTP_ONLY_ACCESS_TOKEN = false
	COOKIE_SECURE_ACCESS_TOKEN    = false

	COOKIE_NAME_REFRESH_TOKEN      = "refresh_token"
	COOKIE_EXPIRES_REFRESH_TOKEN   = 1440  // Minute
	COOKIE_MAX_AGE_REFRESH_TOKEN   = 86400 // Second
	COOKIE_PATH_REFRESH_TOKEN      = "/auth"
	COOKIE_HTTP_ONLY_REFRESH_TOKEN = false
	COOKIE_SECURE_REFRESH_TOKEN    = false

	// Verification Token
	LENGTH_VARIFICATION_TOKEN    = 32
	LIVE_TIME_VARIFICATION_TOKEN = 15 // Minute

	// // Refresh Token
	// LIVE_TIME_REFRESH_TOKEN = 10080 // Minute
	// LENGTH_SALT             = 16
)

// type Conf struct {
// }

// func GetSecretKeyAccessToken() string {
// 	return SECRET_KEY_JWT_ACCESS_TOKEN
// }
// func GetSecretKeyRefreshToken() string {
// 	return SECRET_KEY_JWT_REFRESH_TOKEN
// }

// GetTimeLiveToken() int
// GetHostName() string

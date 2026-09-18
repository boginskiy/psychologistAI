package config

const (

	// Access Token / JWT
	TIME_LIVE_ACCESS_TOKEN = 15 // Minute
	HOST_SITE              = "psychologistAI.com"
	JWT_SECRET_KEY         = "cjlsjdc3r483ucdhcyeruf9ehrc"

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
	COOKIE_PATH_REFRESH_TOKEN      = "/"
	COOKIE_HTTP_ONLY_REFRESH_TOKEN = false
	COOKIE_SECURE_REFRESH_TOKEN    = false

	// Verification Token
	LENGTH_VARIFICATION_TOKEN    = 32
	LIVE_TIME_VARIFICATION_TOKEN = 15 // Minute

	// Refresh Token
	LIVE_TIME_REFRESH_TOKEN = 10080 // Minute
	LENGTH_SALT             = 16
)

package cookie

import "net/http"

type Cooker interface {
	CreateCookie(configName, value string) (*http.Cookie, error)
	ClearCookie(cookie *http.Cookie) (*http.Cookie, error)
	UpdateConfigCookie(Config)
}

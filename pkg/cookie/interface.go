package cookie

import "net/http"

type Cooker interface {
	CreateCookie(configName, value string) (*http.Cookie, error)
	ClearCookie(nameCookie string) (*http.Cookie, error)
	UpdateConfigCookie(Config)
}

package cookie

type Config struct {
	Name     string //
	Expires  int    // Время жизни (Max-Age имеет приоритет над Expires)
	MaxAge   int    // Секунды жизни (86400 = 24 часа). Если -1 — до закрытия браузера. Если 0 — удаление.
	Path     string // Кука будет отправляться на все пути сайта
	HttpOnly bool   // Защита от XSS: JavaScript не сможет прочитать эту куку
	Secure   bool   // Куку можно отправить только по HTTPS
}

func NewConfig(name, path string, expires, maxAge int, httpOnly, secure bool) *Config {
	return &Config{
		Name:     name,
		Expires:  expires,
		MaxAge:   maxAge,
		Path:     path,
		HttpOnly: httpOnly,
		Secure:   secure}
}

var configDefault = Config{
	Name:     "default_token",
	Expires:  5,
	MaxAge:   300,
	Path:     "/",
	HttpOnly: false,
	Secure:   false,
}

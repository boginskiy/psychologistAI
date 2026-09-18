package cookie

import (
	"fmt"
	"net/http"
	"time"
)

var ErrConfigCookie = fmt.Errorf("no config for cookie")

type Cookies struct {
	MapConfig map[string]Config
}

func NewCookies(configs ...Config) *Cookies {
	mapConfig := make(map[string]Config, 10)

	if len(configs) == 0 {
		return &Cookies{MapConfig: mapConfig}
	}

	for _, config := range configs {
		mapConfig[config.Name] = Config{
			Name:     config.Name,
			Path:     config.Path,
			HttpOnly: config.HttpOnly,
			Expires:  config.Expires,
			MaxAge:   config.MaxAge,
			Secure:   config.Secure,
		}
	}
	return &Cookies{MapConfig: mapConfig}
}

// Обновляем config для Cookie
func (c *Cookies) UpdateConfigCookie(config Config) {
	c.MapConfig[config.Name] = Config{
		Name:     config.Name,
		Path:     config.Path,     // Кука будет отправляться на все пути сайта
		HttpOnly: config.HttpOnly, // Защита от XSS: JavaScript не сможет прочитать эту куку
		Secure:   config.Secure,
		Expires:  config.Expires,
		MaxAge:   config.MaxAge,
		// Value:    value,
		// SameSite: http.SameSiteLaxMode,
	}
}

func (c *Cookies) CreateCookie(configName, value string) (*http.Cookie, error) {
	if config, ok := c.MapConfig[configName]; ok {
		return &http.Cookie{
			Name:     config.Name,
			SameSite: http.SameSiteLaxMode,
			Path:     config.Path,
			HttpOnly: config.HttpOnly,
			Secure:   config.Secure,
			Value:    value,
			Expires:  c.transferIntToTime(config.Expires),
			MaxAge:   config.MaxAge,
		}, nil
	}
	return nil, ErrConfigCookie
}

func (c *Cookies) transferIntToTime(tm int) time.Time {
	return time.Now().Add(time.Duration(tm) * time.Minute)
}

// func (c *Cookie) Set(w http.ResponseWriter) {
// 	http.SetCookie(w, c.Cookie)
// }

// func (c *Cookie) Read(r *http.Request) (*http.Cookie, error) {
// 	cookie, err := r.Cookie(c.Name)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return cookie, nil
// }

// func (c *Cookie) ReadAll(r *http.Request) []*http.Cookie {
// 	return r.Cookies()
// }

// func (c *Cookie) Delete(w http.ResponseWriter) {
// 	// Чтобы удалить куку, нужно установить ее заново со значением MaxAge=-1
// 	// или Expires в прошлом времени.
// 	c.MaxAge = -1
// 	c.Expires = time.Unix(0, 0)
// 	http.SetCookie(w, c.Cookie)
// }

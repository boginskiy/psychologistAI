package cooker

import (
	"net/http"
	"time"
)

type CookieManager struct {
	MapConfig map[string]Config
}

func NewCookieManager(configs ...Config) *CookieManager {
	mapConfig := make(map[string]Config, 10)

	if len(configs) == 0 {
		return &CookieManager{MapConfig: mapConfig}
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
	return &CookieManager{MapConfig: mapConfig}
}

// Обновляем config для Cookie
func (c *CookieManager) UpdateConfigCookie(config Config) {
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

func (c *CookieManager) CreateCookie(configName, value string) (*http.Cookie, error) {
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

func (c *CookieManager) ClearCookie(nameCookie string) (*http.Cookie, error) {
	if config, ok := c.MapConfig[nameCookie]; ok {
		return &http.Cookie{
			Name:     config.Name,
			SameSite: http.SameSiteLaxMode,
			Path:     config.Path,
			HttpOnly: config.HttpOnly,
			Secure:   config.Secure,
			Value:    "",              // Пустое значение
			Expires:  time.Unix(0, 1), // Эпоха Unix + 1 секунда (страховка для старых браузеров)
			MaxAge:   -1,              // Удалить немедленно
		}, nil
	}
	return nil, ErrConfigCookie
}

func (c *CookieManager) transferIntToTime(tm int) time.Time {
	return time.Now().Add(time.Duration(tm) * time.Minute)
}

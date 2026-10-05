package utils

import "net/http"

func Scheme(r *http.Request) string {
	// 1. Проверяем заголовок от прокси (стандартный)
	if proto := r.Header.Get("X-Forwarded-Proto"); proto != "" {
		return proto
	}

	// 2. Запасной вариант — другой популярный заголовок
	if proto := r.Header.Get("X-Forwarded-Protocol"); proto != "" {
		return proto
	}

	// 3. Проверяем TLS напрямую
	if r.TLS != nil {
		return "https"
	}

	// 4. По умолчанию — http
	return "http"
}

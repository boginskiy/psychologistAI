package request

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/mileusna/useragent"
)

const (
	CloudProviderHeader = "CF-Connecting-IP"
	XFFHeader           = "X-Forwarded-For"
	ForwardedHeader     = "Forwarded"
)

func ReadAllRequestBody(r *http.Request, item any) (any, error) {
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read request body: %w", err)
	}

	if err := json.Unmarshal(body, &item); err != nil {
		return nil, fmt.Errorf("failed to deserialization request body: %w", err)
	}
	return item, nil
}

// TakeRealUserIP. Извлекает реальный IP клиента, учитывая цепочку прокси.
func TakeRealUserIP(r *http.Request) string {
	// 1. Стандарт от большинства облачных провайдеров и CDN (Cloudflare, Akamai)
	if ip := r.Header.Get(CloudProviderHeader); ip != "" {
		return strings.TrimSpace(ip)
	}

	// 2. Стандарт X-Forwarded-For (XFF)
	// Может содержать цепочку: "клиент, proxy1, proxy2".
	// Нам нужен самый первый (левый) IP.
	if xff := r.Header.Get(XFFHeader); xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}

	// 3. Заголовок Forwarded (RFC 7239)
	// Пример: Forwarded: for=192.0.2.60; proto=https;
	if fwd := r.Header.Get(ForwardedHeader); fwd != "" {
		// Ищем параметр 'for='
		start := strings.Index(strings.ToLower(fwd), "for=")
		if start != -1 {
			val := fwd[start+4:]
			// Обрезаем до первого разделителя (; или ,)
			end := strings.IndexAny(val, ";,")
			if end == -1 {
				end = len(val)
			}

			ip := strings.TrimSpace(val[:end])
			// IPv6 в этом заголовке могут быть в кавычках и квадратных скобках
			ip = strings.TrimPrefix(ip, "\"")
			ip = strings.TrimSuffix(ip, "\"")
			ip = strings.TrimPrefix(ip, "[")
			ip = strings.TrimSuffix(ip, "]")
			return ip
		}
	}

	// 4. RemoteAddr как крайний случай. Содержит IP и порт ("ip:port")
	// Сюда выполнение дойдет, если вы не сидите за прокси (локальная разработка)
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func TakeUserAgent(r *http.Request) string {
	return r.Header.Get("User-Agent")
}

func TakeDeviceInfo(r *http.Request) (os, browser, device string) {
	ua := useragent.Parse(r.Header.Get("User-Agent"))
	return ua.OS, ua.Name, ua.Device
}

package generators

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// GenerateToken создает случайную строку заданной длины в байтах ДО кодирования.
func GenerateToken(byteLength int) (string, error) {
	// Создаем слайс нужного размера
	b := make([]byte, byteLength)

	// Читаем криптографически безопасные случайные числа
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("не удалось прочитать из crypto/rand: %w", err)
	}

	// Кодируем в base64.
	// RawURLEncoding выдает [A-Za-z0-9\-_] и убирает знаки "="
	return base64.RawURLEncoding.EncodeToString(b), nil
}

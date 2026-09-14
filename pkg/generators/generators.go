package generators

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/google/uuid"
)

// GenerateToken создает случайную строку заданной длины в байтах ДО кодирования.
func GenerateTokenBase64(byteLength int) (string, error) {
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

func GenerateRandomBytes(byteLength int) ([]byte, error) {
	// Создаем слайс нужного размера
	b := make([]byte, byteLength)

	// Читаем криптографически безопасные случайные числа
	if _, err := rand.Read(b); err != nil {
		return nil, fmt.Errorf("не удалось прочитать из crypto/rand: %w", err)
	}

	return b, nil
}

func CreateUUIDv7() uuid.UUID {
	return uuid.Must(uuid.NewV7())
}

func CreateUUIDv7ToString() string {
	return uuid.Must(uuid.NewV7()).String()
}

package models

import (
	"time"

	"github.com/boginskiy/psychologistAI/pkg/hashpass"
	"github.com/google/uuid"
)

/*
	AddressIP:
		IP-адрес клиента в момент логина

	UserAgent:
		Полная строка User-Agent браузера. Храните целиком, чтобы можно было проанализировать
		мажорную версию позже.

	CreatedAt:
		Когда была создана запись. Стандартное поле audit log.

	LastUsedAt:
		Очень важное поле. Обновляйте его при каждом успешном использовании refresh-токена. Позволяет чистить
		«мёртвые» сессии (например, удалять те, что не использовались 30 дней)

*/

type DeviceInfo struct {
	OS      string
	Browser string
	Device  string
}

type Session struct {
	ID               string    // jti токена для быстрого поиска
	UserID           uuid.UUID // ID пользователя
	RefreshTokenHash string    // SHA-256 от строки токена

	AddressIP  string     // IP-адрес клиента в момент логина
	UserAgent  string     // Cтрока User-Agent браузера
	DeviceInfo DeviceInfo // Доп инфа

	CreatedAt  time.Time  // Время создания сессии.
	LastUsedAt time.Time  // Время, когда был последний активный сеанс
	RevokedAt  *time.Time // Время отзыва. Если NULL — сессия активна. Если заполнена — токен недействителен.
	ExpiresAt  time.Time  // Время жизни именно этой записи в БД. Должно совпадать с exp внутри Refresh Token.
}

func NewSession(userID uuid.UUID, id, refreshToken, ip, userAgent string, exp time.Time) *Session {
	timeNow := time.Now().UTC()
	return &Session{
		ID:               id,
		UserID:           userID,
		RefreshTokenHash: string(hashpass.CreateBytesHashSHA256(refreshToken)),
		AddressIP:        ip,
		UserAgent:        userAgent,
		// DeviceInfo: TODO...
		CreatedAt:  timeNow,
		LastUsedAt: timeNow,
		ExpiresAt:  exp, // Время берем с refresh token
		RevokedAt:  nil,
	}
}

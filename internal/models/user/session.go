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

type Session struct {
	ID               string    `json:"id" db:"id"`                                 // jti токена для быстрого поиска
	UserID           uuid.UUID `json:"user_id" db:"user_id"`                       // ID пользователя
	RefreshTokenHash []byte    `json:"refresh_token_hash" db:"refresh_token_hash"` // SHA-256 от строки токена

	AddressIP   string `json:"address_ip" db:"address_ip"`     // IP-адрес клиента в момент логина
	UserAgent   string `json:"user_agent" db:"user_agent"`     // Cтрока User-Agent браузера
	OSInfo      string `json:"os_info" db:"os_info"`           // Инфо про OS пользователя
	BrowserInfo string `json:"browser_info" db:"browser_info"` // Инфо про browser пользователя
	DeviceInfo  string `json:"device_info" db:"device_info"`   // Инфо про device пользователя

	PreviousID string `json:"previous_id" db:"previous_id"` // Запоминаем родительскую сессию

	CreatedAt  time.Time  `json:"created_at" db:"created_at"`     // Время создания сессии.
	LastUsedAt time.Time  `json:"last_used_at" db:"last_used_at"` // Время, когда был последний активный сеанс
	RevokedAt  *time.Time `json:"revoked_at" db:"revoked_at"`     // Время отзыва. Если NULL — сессия активна. Если заполнена — токен недействителен.
	ExpiresAt  time.Time  `json:"expires_at" db:"expires_at"`     // Время жизни именно этой записи в БД. Должно совпадать с exp внутри Refresh Token.
}

func NewSession(userID uuid.UUID, id, refreshToken, ip, userAgent string, exp time.Time) *Session {
	timeNow := time.Now().UTC()
	return &Session{
		ID:               id,
		UserID:           userID,
		RefreshTokenHash: hashpass.CreateBytesHashSHA256(refreshToken),
		AddressIP:        ip,
		UserAgent:        userAgent,
		// DeviceInfo: TODO...
		CreatedAt:  timeNow,
		LastUsedAt: timeNow,
		ExpiresAt:  exp, // Время берем с refresh token
		RevokedAt:  nil,
	}
}

package chat

import (
	"time"

	"github.com/boginskiy/psychologistAI/pkg/generators"
	"github.com/google/uuid"
)

type MessageStatus string

const (
	StatusPending   MessageStatus = "pending"   // Ожидание обработки
	StatusComplete  MessageStatus = "complete"  // Успешное завершение
	StatusFailed    MessageStatus = "failed"    // Необработанная ошибка
	StatusStreaming MessageStatus = "streaming" // Потоковая обработка (промежуточный статус)
)

type Message struct {
	ID        uuid.UUID
	ChatID    uuid.UUID // к какому чату относится
	Role      string
	Text      string
	Status    MessageStatus
	CreatedAt *time.Time
	UpdatedAt *time.Time

	// Опционально
	TokenCount int    // сколько токенов в сообщении (для контроля лимитов AI)
	Model      string // какая модель отвечала (gpt-4o, claude-3 и т.д.)
	IsDeleted  bool   // мягкое удаление — пользователь может удалить сообщение
	// Meta       JSONB  // гибкое поле для метаданных (эмоции, теги и т.д.)
}

func NewLiteMessage(text string) *Message {
	// Current time
	timeNow := time.Now().UTC()

	return &Message{
		ID: uuid.UUID([]byte("777")),
		// ChatID:     chatID,
		// Role:       role,
		Text:       text,
		Status:     StatusPending,
		CreatedAt:  &timeNow,
		UpdatedAt:  &timeNow,
		TokenCount: 0,
		Model:      "GPT-3",
		IsDeleted:  false,
	}
}

func NewMessage(chatID uuid.UUID, role, text string) *Message {
	// Current time
	timeNow := time.Now().UTC()

	return &Message{
		ID:         generators.CreateUUIDv7(),
		ChatID:     chatID,
		Role:       role,
		Text:       text,
		Status:     StatusPending,
		CreatedAt:  &timeNow,
		UpdatedAt:  &timeNow,
		TokenCount: 0,
		Model:      "GPT-3",
		IsDeleted:  false,
	}
}

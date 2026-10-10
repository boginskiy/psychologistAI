package chat

import (
	"time"

	"github.com/boginskiy/psychologistAI/internal/adapters/dto"
	"github.com/boginskiy/psychologistAI/pkg/generators"
	"github.com/google/uuid"
)

type Chat struct {
	ID       uuid.UUID
	UserID   uuid.UUID // чей это чат
	Title    string    // название чата (автогенерируется из первого сообщения)
	Messages []Message

	CreatedAt     *time.Time
	UpdatedAt     *time.Time
	LastMessageAt *time.Time // когда последний раз писали — для сортировки списка чатов

	IsActive  bool // активен или архивирован
	IsDeleted bool // мягкое удаление

	// Опционально
	// MessageCount int   // сколько всего сообщений (для пагинации)
	// Meta         JSONB // метаданные
}

func NewChat(infoUser *dto.InfoUser) *Chat {
	// Current time
	timeNow := time.Now().UTC()

	return &Chat{
		ID:     generators.CreateUUIDv7(),
		UserID: infoUser.UserID,

		// TODO. Сделать автогенерацию AI
		Title: makeTitle(infoUser.Message),

		Messages:  make([]Message, 0, 10),
		CreatedAt: &timeNow,
		UpdatedAt: &timeNow,
		// LastMessageAt:
		IsActive:  true,
		IsDeleted: false,
	}
}

func makeTitle(msg string) string {
	if len(msg) > 50 {
		msg = msg[:50] + "_"
	}
	if msg == "" {
		msg = "Новый разговор"
	}
	return msg
}

func (c *Chat) AddMessage(msg ...Message) []Message {
	now := time.Now().UTC()

	c.Messages = append(c.Messages, msg...)

	c.UpdatedAt = &now
	c.LastMessageAt = &now

	return c.Messages
}

func (c *Chat) MessageCount() int {
	return len(c.Messages)
}

func (c *Chat) LastN(n int) []Message {
	if len(c.Messages) <= n {
		return c.Messages
	}
	return c.Messages[len(c.Messages)-n:]
}

package chatrepo

import (
	"fmt"

	"github.com/boginskiy/psychologistAI/internal/db"
	"github.com/boginskiy/psychologistAI/internal/db/mapDB"
	"github.com/boginskiy/psychologistAI/internal/domain/chat"
	"github.com/google/uuid"
)

type ChatRepo struct {
	DB db.DataBase
}

func NewChatRepo() *ChatRepo {
	return &ChatRepo{
		DB: mapDB.NewMapDB(),
	}
}

// func (r *ChatRepo) ExistsBy(userID uuid.UUID) {
// 	tb := r.DB.GetChatTable()
// }

func (r *ChatRepo) Read(userID uuid.UUID) (*chat.Chat, error) {
	tb := r.DB.GetChatTable()

	for _, chat := range tb {
		if chat.UserID == userID {
			return chat, nil
		}
	}
	return nil, fmt.Errorf("no chat for this user")
}

func (r *ChatRepo) Create(chat *chat.Chat) {
	tb := r.DB.GetChatTable()
	tb[chat.ID.String()] = chat
}

func (r *ChatRepo) Update(chat *chat.Chat) {
	tb := r.DB.GetChatTable()
	tb[chat.ID.String()] = chat
}

// TODO
// Сейчас сделано только для одного чата на клиента.
// Функционал многочатовости будет позже

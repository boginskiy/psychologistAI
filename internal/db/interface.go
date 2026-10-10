package db

import (
	"github.com/boginskiy/psychologistAI/internal/domain/chat"
	domain "github.com/boginskiy/psychologistAI/internal/domain/user"
)

type DataBase interface {
	GetUserTable() map[string]*domain.User
	GetSessionTable() map[string]*domain.Session
	GetChatTable() map[string]*chat.Chat
}

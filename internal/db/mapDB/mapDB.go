package mapDB

import (
	chat "github.com/boginskiy/psychologistAI/internal/domain/chat"
	user "github.com/boginskiy/psychologistAI/internal/domain/user"
)

type MapDB struct {
	UserTable    map[string]*user.User    // string == userID
	SessionTable map[string]*user.Session // string == sessionID
	ChatTable    map[string]*chat.Chat    // string == chatID
}

func NewMapDB() *MapDB {
	return &MapDB{
		UserTable:    make(map[string]*user.User, 10),
		SessionTable: make(map[string]*user.Session, 10),
		ChatTable:    make(map[string]*chat.Chat, 10),
	}
}

func (m *MapDB) GetUserTable() map[string]*user.User {
	return m.UserTable
}

func (m *MapDB) GetSessionTable() map[string]*user.Session {
	return m.SessionTable
}

func (m *MapDB) GetChatTable() map[string]*chat.Chat {
	return m.ChatTable
}

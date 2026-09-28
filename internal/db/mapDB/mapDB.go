package mapDB

import (
	domain "github.com/boginskiy/psychologistAI/internal/domain/user"
)

type MapDB struct {
	UserTable    map[string]*domain.User
	SessionTable map[string]*domain.Session
}

func NewMapDB() *MapDB {
	return &MapDB{
		UserTable:    make(map[string]*domain.User, 10),
		SessionTable: make(map[string]*domain.Session, 10),
	}
}

func (m *MapDB) GetUserTable() map[string]*domain.User {
	return m.UserTable
}

func (m *MapDB) GetSessionTable() map[string]*domain.Session {
	return m.SessionTable
}

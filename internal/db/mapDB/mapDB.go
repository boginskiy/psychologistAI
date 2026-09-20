package mapDB

import (
	models "github.com/boginskiy/psychologistAI/internal/models/user"
)

type MapDB struct {
	UserTable    map[string]*models.User
	SessionTable map[string]*models.Session
}

func NewMapDB() *MapDB {
	return &MapDB{
		UserTable:    make(map[string]*models.User, 10),
		SessionTable: make(map[string]*models.Session, 10),
	}
}

func (m *MapDB) GetUserTable() map[string]*models.User {
	return m.UserTable
}

func (m *MapDB) GetSessionTable() map[string]*models.Session {
	return m.SessionTable
}

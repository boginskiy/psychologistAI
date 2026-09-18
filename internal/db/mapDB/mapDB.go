package mapDB

import (
	models "github.com/boginskiy/psychologistAI/internal/models/users"
)

type MapDB struct {
	UserTable map[string]*models.User
}

func NewMapDB() *MapDB {
	return &MapDB{
		UserTable: make(map[string]*models.User, 10),
	}
}

func (m *MapDB) GetUserTable() map[string]*models.User {
	return m.UserTable
}

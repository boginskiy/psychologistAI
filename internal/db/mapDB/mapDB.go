package mapp

import "github.com/boginskiy/psychologistAI/internal/models"

type MapDB struct {
	UserTable map[string]*models.User
}

func (m *MapDB) GetUserTable() map[string]*models.User {
	return m.UserTable
}

package db

import models "github.com/boginskiy/psychologistAI/internal/models/user"

type DataBase interface {
	GetUserTable() map[string]*models.User
	GetSessionTable() map[string]*models.Session
}

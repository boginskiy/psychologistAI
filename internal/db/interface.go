package db

import models "github.com/boginskiy/psychologistAI/internal/models/users"

type DataBase interface {
	GetUserTable() map[string]*models.User
}

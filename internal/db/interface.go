package db

import "github.com/boginskiy/psychologistAI/internal/models"

type DataBase interface {
	GetUserTable() map[string]*models.User
}

package repository

import "github.com/boginskiy/psychologistAI/internal/models"

type UserRepo interface {
	SaveItem(user *models.User) error
}

type CommRepo interface {
}

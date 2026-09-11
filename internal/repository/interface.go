package repository

import "github.com/boginskiy/psychologistAI/internal/models"

type UserRepo interface {
	SaveItem(user *models.User) error
	GetItem(token string) (*models.User, error)
	UpdateItem(user *models.User)
	GetItem2(email string) (*models.User, error)
}

type CommRepo interface {
}

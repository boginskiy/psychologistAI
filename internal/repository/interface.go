package repository

import models "github.com/boginskiy/psychologistAI/internal/models/users"

type UserRepo interface {
	SaveItem(user *models.User) error
	GetItem(token []byte) (*models.User, error)
	UpdateItem(user *models.User)
	GetItem2(email string) (*models.User, error)
}

type CommRepo interface {
}

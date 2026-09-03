package repository

import "github.com/boginskiy/psychologistAI/internal/models"

type UserRepo interface {
	CheckUnic(user *models.User) bool
}

type CommRepo interface {
}

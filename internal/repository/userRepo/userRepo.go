package userrepo

import (
	"fmt"

	"github.com/boginskiy/psychologistAI/internal/db"
	"github.com/boginskiy/psychologistAI/internal/db/mapDB"
	"github.com/boginskiy/psychologistAI/internal/models"
)

type UserRepo struct {
	DB db.DataBase
}

func NewUserRepo() *UserRepo {
	return &UserRepo{
		DB: mapDB.NewMapDB(),
	}
}

func (r *UserRepo) GetItem(token string) (*models.User, error) {
	userTb := r.DB.GetUserTable()
	for _, user := range userTb {
		if user.VerificationToken == token {
			return user, nil
		}
	}
	return nil, fmt.Errorf("user was not found")
}

func (r *UserRepo) UpdateItem(user *models.User) {
	userTb := r.DB.GetUserTable()
	for email := range userTb {
		if email == user.Email {
			userTb[email] = user
			return
		}
	}
}

func (r *UserRepo) SaveItem(user *models.User) error {
	userTb := r.DB.GetUserTable()

	_, ok := userTb[user.Email]
	if ok {
		return fmt.Errorf("email is not unique, try again")
	}

	userTb[user.Email] = user
	return nil
}

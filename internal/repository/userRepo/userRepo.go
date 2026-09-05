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

func (r *UserRepo) SaveItem(user *models.User) error {
	userTb := r.DB.GetUserTable()

	_, ok := userTb[user.Email]
	if ok {
		return fmt.Errorf("user's email is not unique, try again")
	}

	fmt.Printf("%+v", user)

	userTb[user.Email] = user
	return nil
}

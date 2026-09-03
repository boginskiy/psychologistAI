package userrepo

import (
	"github.com/boginskiy/psychologistAI/internal/db"
	"github.com/boginskiy/psychologistAI/internal/models"
)

type UserRepo struct {
	DB db.DataBase
}

func (r *UserRepo) CheckUnic(user *models.User) bool {
	userTb := r.DB.GetUserTable()

	_, ok := userTb[user.Email]
	if ok {
		return !ok
	}
	return ok
}

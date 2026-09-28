package userrepo

import (
	"fmt"

	"github.com/boginskiy/psychologistAI/internal/db"
	"github.com/boginskiy/psychologistAI/internal/db/mapDB"

	domain "github.com/boginskiy/psychologistAI/internal/domain/user"
)

type UserRepo struct {
	DB db.DataBase
}

func NewUserRepo() *UserRepo {
	return &UserRepo{
		DB: mapDB.NewMapDB(),
	}
}

func (r *UserRepo) ReadByEmail(email string) (*domain.User, error) {
	userTb := r.DB.GetUserTable()
	for e, user := range userTb {
		if e == email {
			return user, nil
		}
	}
	return nil, fmt.Errorf("user was not found")
}

func (r *UserRepo) ReadByToken(hashToken []byte) (*domain.User, error) {
	userTb := r.DB.GetUserTable()
	for _, user := range userTb {

		if user.CompareHash(hashToken) {
			return user, nil
		}
	}
	return nil, fmt.Errorf("user was not found")
}

func (r *UserRepo) UpdateItem(user *domain.User) {
	userTb := r.DB.GetUserTable()
	for email := range userTb {
		if email == user.Email {
			userTb[email] = user
			return
		}
	}
}

func (r *UserRepo) Create(user *domain.User) error {
	userTb := r.DB.GetUserTable()

	_, ok := userTb[user.Email]
	if ok {
		return fmt.Errorf("email is not unique, try again")
	}

	userTb[user.Email] = user
	return nil
}

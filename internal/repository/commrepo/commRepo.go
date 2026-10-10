package commrepo

import (
	"github.com/boginskiy/psychologistAI/internal/db"
	domain "github.com/boginskiy/psychologistAI/internal/domain/user"
)

type CommRepo struct {
	DB db.DataBase
}

func (r *CommRepo) CheckUnic(user *domain.User) bool {

	return false
}

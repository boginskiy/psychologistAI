package commrepo

import (
	"github.com/boginskiy/psychologistAI/internal/db"
	models "github.com/boginskiy/psychologistAI/internal/models/user"
)

type CommRepo struct {
	DB db.DataBase
}

func (r *CommRepo) CheckUnic(user *models.User) bool {

	return false
}

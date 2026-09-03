package commrepo

import (
	"github.com/boginskiy/psychologistAI/internal/db"
	"github.com/boginskiy/psychologistAI/internal/models"
)

type CommRepo struct {
	DB db.DataBase
}

func (r *CommRepo) CheckUnic(user *models.User) bool {

	return false
}

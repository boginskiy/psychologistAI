package userrepo

import (
	"github.com/boginskiy/psychologistAI/internal/db"
	"github.com/boginskiy/psychologistAI/internal/db/mapDB"
	models "github.com/boginskiy/psychologistAI/internal/models/user"
)

type SessionRepo struct {
	DB db.DataBase
}

func NewSessionRepo() *SessionRepo {
	return &SessionRepo{
		DB: mapDB.NewMapDB(),
	}
}

func (r *SessionRepo) Create(session *models.Session) error {
	tb := r.DB.GetSessionTable()
	tb[session.ID] = session
	return nil
}

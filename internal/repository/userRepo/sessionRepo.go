package userrepo

import (
	"fmt"
	"time"

	"github.com/boginskiy/psychologistAI/internal/db"
	"github.com/boginskiy/psychologistAI/internal/db/mapDB"
	models "github.com/boginskiy/psychologistAI/internal/models/user"
	"github.com/google/uuid"
)

type SessionRepo struct {
	DB db.DataBase
}

func NewSessionRepo() *SessionRepo {
	return &SessionRepo{
		DB: mapDB.NewMapDB(),
	}
}

func (r *SessionRepo) CancelSessions(userID uuid.UUID) {
	tb := r.DB.GetSessionTable()
	for _, session := range tb {
		if session.UserID == userID {
			timeNow := time.Now().UTC()
			session.RevokedAt = &timeNow
		}
	}
}

func (r *SessionRepo) CancelSession(sessionID string) {
	tb := r.DB.GetSessionTable()
	if session, ok := tb[sessionID]; ok {
		timeNow := time.Now().UTC()
		session.RevokedAt = &timeNow
	}
}

func (r *SessionRepo) Read(id string) (*models.Session, error) {
	tb := r.DB.GetSessionTable()
	for key, session := range tb {
		if key == id {
			return session, nil
		}
	}
	return nil, fmt.Errorf("there is no session")
}

func (r *SessionRepo) Create(session *models.Session) error {
	tb := r.DB.GetSessionTable()
	tb[session.ID] = session
	return nil
}

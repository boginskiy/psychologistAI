package repository

import (
	models "github.com/boginskiy/psychologistAI/internal/models/user"
	"github.com/google/uuid"
)

type UserReader interface {
	ReadByToken(token []byte) (*models.User, error)
	ReadByEmail(email string) (*models.User, error)
}

type UserCreater interface {
	Create(user *models.User) error
}

type UserUpdater interface {
	UpdateItem(user *models.User)
}

type UserRepo interface {
	UserUpdater
	UserCreater
	UserReader
}

// =============================
type SessionReader interface {
	Read(id string) (*models.Session, error)
}

type SessionCreater interface {
	Create(session *models.Session) error
}

type SessionRepo interface {
	// SessionUpdater
	SessionReader
	SessionCreater

	CancelSessions(userID uuid.UUID)
	CancelSession(sessionID string)
}

type CommRepo interface {
}

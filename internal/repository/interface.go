package repository

import (
	domain "github.com/boginskiy/psychologistAI/internal/domain/user"
	"github.com/google/uuid"
)

type UserReader interface {
	ReadByToken(token []byte) (*domain.User, error)
	ReadByEmail(email string) (*domain.User, error)
}

type UserCreater interface {
	Create(user *domain.User) error
}

type UserUpdater interface {
	UpdateItem(user *domain.User)
}

type UserRepo interface {
	UserUpdater
	UserCreater
	UserReader
}

// =============================
type SessionReader interface {
	Read(id string) (*domain.Session, error)
}

type SessionCreater interface {
	Create(session *domain.Session) error
}

type SessionRepo interface {
	// SessionUpdater
	SessionReader
	SessionCreater

	UpdateAfterRefresh(newSession *domain.Session, offset int) error
	DeleteSession(sessionID string, offset int)
	IsActiveSession(sessionID string) bool
	CancelSessions(userID uuid.UUID)
	CancelSession(sessionID string)
}

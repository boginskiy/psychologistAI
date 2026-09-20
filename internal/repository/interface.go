package repository

import models "github.com/boginskiy/psychologistAI/internal/models/user"

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
type SessionCreater interface {
	Create(session *models.Session) error
}

type SessionRepo interface {
	// SessionUpdater
	// SessionReader
	SessionCreater
}

type CommRepo interface {
}

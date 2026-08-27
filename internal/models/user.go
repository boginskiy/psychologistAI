package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	// Основные поля
	ID       uuid.UUID `json:"id" db:"id"`
	Email    string    `json:"email" db:"email" validate:"required,email"`
	Password string    `json:"-" db:"password_hash"` // Храним хеш, не выводим

	// Личная информация
	FirstName string `json:"first_name,omitempty" db:"first_name"`
	LastName  string `json:"last_name,omitempty" db:"last_name"`
	Phone     string `json:"phone,omitempty" db:"phone"`

	// Статусы
	IsActive bool   `json:"is_active" db:"is_active"`
	IsAdmin  bool   `json:"is_admin" db:"is_admin"`
	Role     string `json:"role" db:"role"` // user, admin, moderator

	// Временные метки
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty" db:"last_login_at"`
	DeletedAt   *time.Time `json:"-" db:"deleted_at"` // Soft delete
}

package dto

import "github.com/google/uuid"

// CreateUser - DTO для создания пользователя
type CreateUser struct {
	Email            string `json:"email" db:"email" validate:"required,email"`
	Password         string `json:"password" db:"password_hash"`
	Name             string `json:"name,omitempty" db:"name"`
	Phone            string `json:"phone,omitempty" db:"phone"`
	VerificationLink string `json:"-" db:"-"`
}

// LoginUser - DTO для авторизации пользователя
type LoginUser struct {
	Email     string `json:"email" db:"email" validate:"required,email"`
	Password  string `json:"password" db:"password_hash"`
	IP        string `json:"-" db:"user_ip"`
	UserAgent string `json:"-" db:"user_agent"`
	OS        string
	Browser   string
	Device    string
}

// UpdateUser - DTO для обновления пользователя
type UpdateUser struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
	Phone string `json:"phone,omitempty" db:"phone"`
}

// InfoUser - DTO информация о пользователе
type InfoUser struct {
	UserID    uuid.UUID `json:"uid" db:"uid"`
	UserRoles []string  `json:"roles" db:"roles"`
	TokenType string    `json:"type" db:"type"`
	SessionID string    `json:"session_id" db:"session_id"`
}

package service

import (
	"context"

	"github.com/boginskiy/psychologistAI/internal/adapters/dto"
	"github.com/boginskiy/psychologistAI/internal/models"
)

type UserService struct {
	Validater Validater
}

func NewUserService(validater Validater) *UserService {
	return &UserService{
		Validater: validater,
	}
}

func (s *UserService) CreateUser(ctx context.Context, userReq *dto.CreateUserRequest) (*dto.UserResponse, error) {
	// Валидация Email
	err := s.Validater.CheckNotEmptyStrField("email", userReq.Email)
	if err != nil {
		return nil, err
	}

	// Валидация Password
	err = s.Validater.CheckNotEmptyStrField("password", userReq.Password)
	if err != nil {
		return nil, err
	}

	// Domain user
	newUser := models.NewUser(userReq)

	return nil, nil
}

// type User struct {
// 	// Основные поля
// 	ID       uuid.UUID `json:"id" db:"id"`
// 	Email    string    `json:"email" db:"email" validate:"required,email"`
// 	Password string    `json:"-" db:"password_hash"` // Храним хеш, не выводим

// 	// Личная информация
// 	FirstName string `json:"first_name,omitempty" db:"first_name"`
// 	LastName  string `json:"last_name,omitempty" db:"last_name"`
// 	Phone     string `json:"phone,omitempty" db:"phone"`

// 	// Статусы
// 	IsActive bool   `json:"is_active" db:"is_active"`
// 	IsAdmin  bool   `json:"is_admin" db:"is_admin"`
// 	Role     string `json:"role" db:"role"` // user, admin, moderator

// 	// Временные метки
// 	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
// 	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
// 	LastLoginAt    *time.Time `json:"last_login_at,omitempty" db:"last_login_at"`
// 	LastActivityAt *time.Time `json:"last_activity_at,omitempty" db:"last_activity_at"`
// 	DeletedAt      *time.Time `json:"-" db:"deleted_at"` // Soft delete
// }

// type CreateUserRequest struct {
// 	ID        uuid.UUID `json:"id" db:"id"`
// 	Email     string    `json:"email" db:"email" validate:"required,email"`
// 	Password  string    `json:"-" db:"password_hash"`
// 	FirstName string    `json:"first_name,omitempty" db:"first_name"`
// 	LastName  string    `json:"last_name,omitempty" db:"last_name"`
// 	Phone     string    `json:"phone,omitempty" db:"phone"`
// }

package service

import (
	"context"
	"fmt"
)

type ValidService struct {
}

func NewValidService(ctx context.Context) *ValidService {
	return &ValidService{}
}

func (s *ValidService) CheckNotEmptyStrField(name, field string) error {
	if field == "" {
		return fmt.Errorf("field '%s' cannot be empty string", name)
	}
	return nil
}

// type CreateUserRequest struct {
// 	ID        uuid.UUID `json:"id" db:"id"`
// 	Email     string    `json:"email" db:"email" validate:"required,email"`
// 	Password  string    `json:"-" db:"password_hash"`
// 	FirstName string    `json:"first_name,omitempty" db:"first_name"`
// 	LastName  string    `json:"last_name,omitempty" db:"last_name"`
// 	Phone     string    `json:"phone,omitempty" db:"phone"`
// }

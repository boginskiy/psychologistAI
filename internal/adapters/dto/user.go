package dto

// CreateUserRequest - DTO для создания пользователя
type CreateUserRequest struct {
	Email     string `json:"email" db:"email" validate:"required,email"`
	Password  string `json:"password" db:"password_hash"`
	FirstName string `json:"first_name,omitempty" db:"first_name"`
	LastName  string `json:"last_name,omitempty" db:"last_name"`
	Phone     string `json:"phone,omitempty" db:"phone"`
}

// UpdateUserRequest - DTO для обновления пользователя
type UpdateUserRequest struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Email     string `json:"email,omitempty"`
	Phone     string `json:"phone,omitempty" db:"phone"`
}

// UserResponse - DTO для ответа после создания пользователя
// type UserResponse struct {
// 	Email     string    `json:"email"`
// 	FirstName string    `json:"first_name,omitempty"`
// 	LastName  string    `json:"last_name,omitempty"`
// 	CreatedAt time.Time `json:"created_at"`
// }

package dto

// CreateUserRequest - DTO для создания пользователя
type CreateUser struct {
	Email    string `json:"email" db:"email" validate:"required,email"`
	Password string `json:"password" db:"password_hash"`
	Name     string `json:"name,omitempty" db:"name"`
	Phone    string `json:"phone,omitempty" db:"phone"`
}

// LoginUserRequest - DTO для авторизации пользователя
type LoginUser struct {
	Email    string `json:"email" db:"email" validate:"required,email"`
	Password string `json:"password" db:"password_hash"`
}

// UpdateUserRequest - DTO для обновления пользователя
type UpdateUser struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Email     string `json:"email,omitempty"`
	Phone     string `json:"phone,omitempty" db:"phone"`
}

package response

import (
	"time"

	"github.com/boginskiy/psychologistAI/internal/api/errors"
)

type UserResponse struct {
	Name      string     `json:"name,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	errors.APIError
}

func NewUserResponse() *UserResponse {
	return &UserResponse{}
}

func (r *UserResponse) GetStatus() int {
	return r.Status
}

func (r *UserResponse) PrepareOKResponse(status int) {
	r.Status = status
}

func (r *UserResponse) PrepareErrResponse(err error, status int) {
	r.Message = err.Error()
	r.Status = status
}

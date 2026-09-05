package adapters

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/boginskiy/psychologistAI/internal/adapters/dto"
	"github.com/boginskiy/psychologistAI/internal/api/response"
	"github.com/boginskiy/psychologistAI/internal/models"
)

func ToCreateUserRequest(r *http.Request) (*dto.CreateUserRequest, error) {
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}

	var user dto.CreateUserRequest
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	return &user, nil
}

func ToUserResponse(user *models.User) *response.UserResponse {
	return &response.UserResponse{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: &user.CreatedAt,
	}
}

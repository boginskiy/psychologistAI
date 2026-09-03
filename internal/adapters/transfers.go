package adapters

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/boginskiy/psychologistAI/internal/adapters/dto"
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

func ToUserResponse(user *models.User) *dto.UserResponse {
	return &dto.UserResponse{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: user.CreatedAt,
	}
}

func ToCurrentTime(r *http.Request, tm time.Time) time.Time {
	// Take Tz
	tz := r.Header.Get("Accept-Timezone")
	if tz == "" {
		return tm
	}
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return tm
	}
	return tm.In(loc)
}

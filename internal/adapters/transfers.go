package adapters

import (
	"net/http"

	"github.com/boginskiy/psychologistAI/internal/api/response"
	"github.com/boginskiy/psychologistAI/internal/user/models"
	"github.com/go-chi/chi"
)

func ToLoginUserRequest(r *http.Request) {
	// TODO/ Чтение боди, адаптеры обновить. Перейти на логин и далее делать логин
}

func ToToken(r *http.Request) string {
	return chi.URLParam(r, "token")
}

func ToUserResponse(user *models.User) *response.UserResponse {
	return &response.UserResponse{
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
	}
}

func ToUserResponseOnlyMess(msg string) *response.UserResponse {
	tmp := &response.UserResponse{}
	tmp.Message = msg
	return tmp
}

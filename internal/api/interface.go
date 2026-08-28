package api

import (
	"net/http"

	"github.com/boginskiy/psychologistAI/internal/adapters/dto"
)

type Response interface {
	SendResponse(w http.ResponseWriter, userRes *dto.UserResponse, status int)
	SendError(w http.ResponseWriter, err error, status int)
}

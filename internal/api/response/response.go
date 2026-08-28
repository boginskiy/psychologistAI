package response

import (
	"encoding/json"
	"net/http"

	"github.com/boginskiy/psychologistAI/internal/adapters/dto"
	"github.com/boginskiy/psychologistAI/internal/api/errors"
)

type Resp struct {
}

func (r *Resp) SendError(w http.ResponseWriter, err error, status int) {
	apiErr := errors.NewAPIError(err, status)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(apiErr)
}

func (r *Resp) SendResponse(w http.ResponseWriter, userRes *dto.UserResponse, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(userRes)
}

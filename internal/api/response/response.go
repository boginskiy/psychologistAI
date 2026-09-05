package response

import (
	"encoding/json"
	"net/http"

	"github.com/boginskiy/psychologistAI/internal/api"
)

type Response struct {
}

func NewResponse() *Response {
	return &Response{}
}

func (r *Response) SendResponse(w http.ResponseWriter, item api.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(item.GetStatus())
	json.NewEncoder(w).Encode(item)
}

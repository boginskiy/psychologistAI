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

func (r *Response) SendResponse(w http.ResponseWriter, body api.ResponseBody) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(body.GetStatus())
	json.NewEncoder(w).Encode(body)
}

func (r *Response) AddSetCookies(w http.ResponseWriter, cookies ...*http.Cookie) {
	if len(cookies) == 0 {
		return
	}
	for _, cookie := range cookies {
		http.SetCookie(w, cookie)
	}
}

func (r *Response) SetContentType(w http.ResponseWriter, contentType string) {
	w.Header().Set("Content-Type", contentType)
}

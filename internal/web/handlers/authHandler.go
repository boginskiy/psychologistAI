package handlers

import (
	"net/http"

	"github.com/boginskiy/psychologistAI/internal/middleware"
	"github.com/go-chi/chi"
)

type AuthHandler struct {
	basepath string
}

func NewAuthHandler(
	bpath string,
) *AuthHandler {

	return &AuthHandler{
		basepath: bpath,
	}
}

func (h *AuthHandler) Registration(r chi.Router, middleware middleware.HandleMiddleware) {
	r.Route(h.basepath, func(r chi.Router) {
		r.Post("/login", h.Loginer)

	})
}

func (h *AuthHandler) Loginer(w http.ResponseWriter, r *http.Request) {
	// TODO. Переход на LOG IN, но так чтобы не дургать ручку, как это сделать ?
}

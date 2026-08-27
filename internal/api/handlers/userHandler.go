package handlers

import (
	"net/http"

	"github.com/boginskiy/psychologistAI/internal/service"
	"github.com/go-chi/chi"
)

type UserHandler struct {
	UserService service.UserService
}

func (h *UserHandler) Registration(r chi.Router) {
	r.Get("/", listUsers)           // GET /api/v1/user
	r.Post("/register", h.Register) // POST /api/v1/user/register

	r.Route("/{userID}", func(r chi.Router) {
		r.Get("/", getUser)    // GET /api/v1/user/{userID}
		r.Put("/", updateUser) // PUT /api/v1/user/{userID}

	})
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	// ID       uuid.UUID  Генерируем v7 (упорядоченный)

	h.UserService.CreateNewUser(r.Context())

}

// r.Get("/", h.List)          // GET /api/v1/user
//     r.Post("/", h.Create)       // POST /api/v1/user
//     r.Get("/{id}", h.Get)       // GET /api/v1/user/{id}
//     r.Put("/{id}", h.Update)    // PUT /api/v1/user/{id}
//     r.Delete("/{id}", h.Delete) // DELETE /api/v1/user/{id}

// TODO
// Регистрация

package handlers

import (
	"net/http"

	"github.com/boginskiy/psychologistAI/internal/adapters"
	"github.com/boginskiy/psychologistAI/internal/api"
	"github.com/boginskiy/psychologistAI/internal/service"
	"github.com/go-chi/chi"
)

type UserHandler struct {
	UserService service.UserService
	Response    api.Response
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
	// Adapter
	userReq, err := adapters.ToCreateUserRequest(r)
	if err != nil {
		h.Response.SendError(w, err, http.StatusBadRequest)
		return
	}

	// Service
	userRes, err := h.UserService.CreateUser(r.Context(), userReq)

	// Errors

	// Response
	h.Response.SendResponse(w, userRes, http.StatusOK)
}

// TODO
//

// ID       uuid.UUID  Генерируем v7 (упорядоченный)
// r.Get("/", h.List)          // GET /api/v1/user
//     r.Post("/", h.Create)       // POST /api/v1/user
//     r.Get("/{id}", h.Get)       // GET /api/v1/user/{id}
//     r.Put("/{id}", h.Update)    // PUT /api/v1/user/{id}
//     r.Delete("/{id}", h.Delete) // DELETE /api/v1/user/{id}

// Сущность:
// UserRequest, UserResponse
// UserDB, UserModel

// TODO
// Регистрация
// Авторизация
//

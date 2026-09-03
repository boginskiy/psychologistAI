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
	basepath    string
}

func NewUserHandler(bpath string, userServ service.UserService, resp api.Response) *UserHandler {
	return &UserHandler{
		UserService: userServ,
		Response:    resp,
		basepath:    bpath,
	}
}

func (h *UserHandler) Registration(r chi.Router) {
	r.Route(h.basepath, func(r chi.Router) {
		r.Post("/register", h.Register) // POST /api/v1/user/register
	})

	// r.Get("/", listUsers)           // GET /api/v1/user
	// r.Route("/{userID}", func(r chi.Router) {
	// 	r.Get("/", getUser)    // GET /api/v1/user/{userID}
	// 	r.Put("/", updateUser) // PUT /api/v1/user/{userID}

	// })
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	// Adapters
	userReq, err := adapters.ToCreateUserRequest(r)
	if err != nil {
		h.Response.SendError(w, err, http.StatusBadRequest)
		return
	}

	// Service
	userRes, err := h.UserService.CreateUser(r.Context(), userReq)

	// Errors
	if err != nil {
		h.Response.SendError(w, err, http.StatusBadRequest)
		return
	}

	// Обработка ошибок, отдельная история
	// Пока все BadRequest

	// Current Time
	userRes.CreatedAt = adapters.ToCurrentTime(r, userRes.CreatedAt)

	// Response
	h.Response.SendResponse(w, userRes, http.StatusOK)
}

// // Регистрируем маршруты.
// r.Route("/users", userHandler.Registration)
// r.Route("/home", homeHandler.Registration)

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

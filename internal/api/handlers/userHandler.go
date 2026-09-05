package handlers

import (
	"net/http"

	"github.com/boginskiy/psychologistAI/internal/adapters"
	"github.com/boginskiy/psychologistAI/internal/api"
	"github.com/boginskiy/psychologistAI/internal/api/response"
	"github.com/boginskiy/psychologistAI/internal/service"
	"github.com/boginskiy/psychologistAI/pkg/timeproc"
	"github.com/go-chi/chi"
)

type UserHandler struct {
	UserService service.UserService
	Sender      api.Sender
	basepath    string
}

func NewUserHandler(bpath string, userServ service.UserService, sender api.Sender) *UserHandler {
	return &UserHandler{
		UserService: userServ,
		Sender:      sender,
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

// func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
// 	state := "random-state-string" // обязательно для защиты от CSRF
// 	http.Redirect(w, r, oauth2Config.AuthCodeURL(state), http.StatusFound)
// }

// // Эндпоинт callback — обмен кода на токены и получение user info
// func handleCallback(w http.ResponseWriter, r *http.Request) {
// 	ctx := context.Background()

// 	// Проверяем state (опущено для краткости)

// 	// 1. Обмениваем code на токены
// 	oauth2Token, err := oauth2Config.Exchange(ctx, r.URL.Query().Get("code"))
// 	if err != nil {
// 		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
// 		return
// 	}

// 	// 2. Извлекаем ID Token (JWT) из ответа
// 	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
// 	if !ok {
// 		http.Error(w, "Missing ID Token", http.StatusInternalServerError)
// 		return
// 	}

// 	// 3. Верифицируем ID Token
// 	idToken, err := idTokenVerifier.Verify(ctx, rawIDToken)
// 	if err != nil {
// 		http.Error(w, "Failed to verify ID Token", http.StatusInternalServerError)
// 		return
// 	}

// 	// 4. Извлекаем claims (информацию о пользователе)
// 	var claims struct {
// 		Email         string `json:"email"`
// 		EmailVerified bool   `json:"email_verified"`
// 		Name          string `json:"name"`
// 		Picture       string `json:"picture"`
// 		Sub           string `json:"sub"` // уникальный идентификатор пользователя
// 	}
// 	if err := idToken.Claims(&claims); err != nil {
// 		http.Error(w, "Failed to parse claims", http.StatusInternalServerError)
// 		return
// 	}

// 	// 5. Работаем с пользовательскими данными
// 	log.Printf("User: %s (%s)", claims.Name, claims.Email)

// 	// Здесь можно создать сессию, сохранить пользователя в БД и т.д.
// 	w.Write([]byte("Hello, " + claims.Name))
// }

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	userResponse := response.NewUserResponse()

	// Adapters
	userRequest, err := adapters.ToCreateUserRequest(r)
	if err != nil {
		userResponse.PrepareErrResponse(err, http.StatusBadRequest)
		h.Sender.SendResponse(w, userResponse, http.StatusBadRequest)
		return
	}

	// Service
	userTmp, err := h.UserService.CreateUser(r.Context(), userRequest)

	// Errors
	if err != nil {
		userResponse.PrepareErrResponse(err, http.StatusBadRequest)
		h.Sender.SendResponse(w, userResponse, http.StatusBadRequest)
		return
	}

	// Update Time
	userTmp.CreatedAt = timeproc.ConvertTimeUtcToLocalByRequest(userTmp.CreatedAt, r)

	userResponse = userTmp
	userResponse.PrepareOKResponse(http.StatusOK)
	h.Sender.SendResponse(w, userResponse, http.StatusOK)
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

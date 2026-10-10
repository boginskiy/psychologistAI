package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/boginskiy/psychologistAI/cmd/config"
	"github.com/boginskiy/psychologistAI/internal/api"
	"github.com/boginskiy/psychologistAI/internal/api/adapters"
	"github.com/boginskiy/psychologistAI/internal/errs"
	"github.com/boginskiy/psychologistAI/internal/middleware"
	"github.com/boginskiy/psychologistAI/internal/msgs"

	"github.com/boginskiy/psychologistAI/internal/api/response"
	"github.com/boginskiy/psychologistAI/internal/service"

	"github.com/go-chi/chi"
)

type UserHandler struct {
	UserService service.UserService
	Responder   api.Responder
	Requester   api.Requester
	basepath    string
}

func NewUserHandler(bpath string, userServ service.UserService, responder api.Responder, requester api.Requester) *UserHandler {
	return &UserHandler{
		basepath:    bpath,
		UserService: userServ,
		Responder:   responder,
		Requester:   requester,
	}
}

func (h *UserHandler) Registration(r chi.Router, middleware middleware.HandleMiddleware) {
	r.Route(h.basepath, func(r chi.Router) {
		// Public
		r.Group(func(r chi.Router) {
			r.Post("/registration", h.Register)        // POST /api/v1/user/registration
			r.Get("/verification/{token}", h.Verifier) // GET  /api/v1/user/verification/{token}
		})

		// TODO...
		// r.Get("/{id}", h.Informer) // GET  /api/v1/user/{id}
		// /api/profile

		// TODO ...
		// Имя пользователя фронтенду обычно не нужно для вызова эндпоинта
		// /api/get-client-notes. Если имя понадобится для отображения в UI,
		// фронт сделает один отдельный запрос /api/profile после успешного входа,
		// получит полные данные и закэширует их у себя во Vuex/Redux.
		// Засорять каждый HTTP-запрос именем неэффективно.

	})
}

func (h *UserHandler) Verifier(w http.ResponseWriter, r *http.Request) {
	body := &response.InfoBody{}
	token := chi.URLParam(r, "token")

	// Service
	_, err := h.UserService.Verification(r.Context(), token)

	// Errors
	if err != nil {
		log.Printf("error: %v\n", err) // + logger

		switch {
		// Verification
		case errors.Is(err, errs.ErrLinkVerification):
			body.ErrorUpdate(errs.ErrLinkVerification, http.StatusNotFound)
		case errors.Is(err, errs.ErrRepeatVerification):
			body.ErrorUpdate(errs.ErrRepeatVerification, http.StatusBadRequest)
		case errors.Is(err, errs.ErrAttemptsVerification):
			body.ErrorUpdate(errs.ErrAttemptsVerification, http.StatusTooManyRequests)

		// Server
		case errors.Is(err, errs.ErrServer):
			body.ErrorUpdate(errs.ErrServer, http.StatusInternalServerError)

		default:
			body.ErrorUpdate(err, http.StatusBadRequest)
		}
		h.Responder.SendResponse(w, body)
		return
	}

	body.InfoUpdate(msgs.MessOkVerification, http.StatusOK)
	h.Responder.SendResponse(w, body)
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	createUser, err := adapters.ToCreateUserFromRequest(r)
	body := &response.InfoBody{}

	if err != nil {
		body.ErrorUpdate(err, http.StatusBadRequest)
		h.Responder.SendResponse(w, body)
		return
	}

	// Service
	userDomen, err := h.UserService.Create(r.Context(), createUser)

	// Errors
	if err != nil {
		switch {
		// Credentials
		case errors.Is(err, errs.ErrInvalidCredentials):
			body.ErrorUpdate(errs.ErrInvalidCredentials, http.StatusUnauthorized)

			// Server
		case errors.Is(err, errs.ErrServer):
			body.ErrorUpdate(errs.ErrServer, http.StatusInternalServerError)

		default:
			body.ErrorUpdate(err, http.StatusBadRequest)
		}
		h.Responder.SendResponse(w, body)
		return
	}

	msg := msgs.FuncNeedRegistration(userDomen.Email, config.LIVE_TIME_VARIFICATION_TOKEN)
	body.InfoUpdate(msg, http.StatusOK)
	h.Responder.SendResponse(w, body)
}

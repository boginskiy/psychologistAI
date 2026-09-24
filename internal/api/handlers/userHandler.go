package handlers

import (
	"errors"
	"net/http"

	"github.com/boginskiy/psychologistAI/cmd/config"
	"github.com/boginskiy/psychologistAI/internal/adapters/dto"
	"github.com/boginskiy/psychologistAI/internal/api"
	"github.com/boginskiy/psychologistAI/internal/api/vars"
	"github.com/boginskiy/psychologistAI/internal/errs"

	"github.com/boginskiy/psychologistAI/internal/models/response"
	"github.com/boginskiy/psychologistAI/internal/service"

	"github.com/boginskiy/psychologistAI/internal/api/request"
	"github.com/go-chi/chi"
)

const Token = "token"

type UserHandler struct {
	UserService    service.UserService
	ResponseSender api.ResponseSender
	basepath       string
}

func NewUserHandler(bpath string, userServ service.UserService, resSender api.ResponseSender) *UserHandler {
	return &UserHandler{
		UserService:    userServ,
		ResponseSender: resSender,
		basepath:       bpath,
	}
}

func (h *UserHandler) Registration(r chi.Router) {
	r.Route(h.basepath, func(r chi.Router) {
		r.Post("/registration", h.Register)        // POST /api/v1/user/registration
		r.Get("/verification/{token}", h.Verifier) // GET  /api/v1/user/verification/{token}

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
	infoResponse := &response.InfoResponse{}
	token := chi.URLParam(r, Token)

	_, err := h.UserService.Verification(r.Context(), token)

	// Errors
	if err != nil {
		switch {
		// Verification
		case errors.Is(err, errs.ErrLinkVerification):
			infoResponse.ErrorUpdate(errs.ErrLinkVerification, http.StatusNotFound)
		case errors.Is(err, errs.ErrRepeatVerification):
			infoResponse.ErrorUpdate(errs.ErrRepeatVerification, http.StatusBadRequest)
		case errors.Is(err, errs.ErrAttemptsVerification):
			infoResponse.ErrorUpdate(errs.ErrAttemptsVerification, http.StatusTooManyRequests)

		// Server
		case errors.Is(err, errs.ErrServer):
			infoResponse.ErrorUpdate(errs.ErrServer, http.StatusInternalServerError)

		default:
			infoResponse.ErrorUpdate(err, http.StatusBadRequest)
		}
		h.ResponseSender.SendResponse(w, infoResponse)
		return
	}

	infoResponse.InfoUpdate(vars.MessOkVerification, http.StatusOK)
	h.ResponseSender.SendResponse(w, infoResponse)
}

// Убрать из сервиса подготовку user Response и перенести ее сюда
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	infoResponse := &response.InfoResponse{}
	createUser := &dto.CreateUser{}

	_, err := request.ReadAllRequestBody(r, createUser)

	if err != nil {
		infoResponse.ErrorUpdate(err, http.StatusBadRequest)
		h.ResponseSender.SendResponse(w, infoResponse)
		return
	}

	// Service
	userDomen, err := h.UserService.Create(r.Context(), createUser)

	// Errors
	if err != nil {
		switch {
		// Credentials
		case errors.Is(err, errs.ErrInvalidCredentials):
			infoResponse.ErrorUpdate(errs.ErrInvalidCredentials, http.StatusUnauthorized)

			// Server
		case errors.Is(err, errs.ErrServer):
			infoResponse.ErrorUpdate(errs.ErrServer, http.StatusInternalServerError)

		default:
			infoResponse.ErrorUpdate(err, http.StatusBadRequest)
		}
		h.ResponseSender.SendResponse(w, infoResponse)
		return
	}

	msg := vars.FuncNeedRegistration(userDomen.Email, config.LIVE_TIME_VARIFICATION_TOKEN)
	infoResponse.InfoUpdate(msg, http.StatusOK)
	h.ResponseSender.SendResponse(w, infoResponse)
}

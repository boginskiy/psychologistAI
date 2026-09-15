package handlers

import (
	"errors"
	"net/http"

	"github.com/boginskiy/psychologistAI/internal/adapters/dto"
	"github.com/boginskiy/psychologistAI/internal/api"
	"github.com/boginskiy/psychologistAI/internal/api/vars"
	"github.com/boginskiy/psychologistAI/internal/errs/users"
	"github.com/boginskiy/psychologistAI/internal/models/responses"
	models "github.com/boginskiy/psychologistAI/internal/models/users"
	"github.com/boginskiy/psychologistAI/internal/service"

	"github.com/boginskiy/psychologistAI/pkg/request"
	"github.com/go-chi/chi"
)

const Token = "token"

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
		r.Post("/registration", h.Register) // POST /api/v1/user/registration
		r.Post("/login", h.Loginer)         // POST /api/v1/user/login

		r.Get("/verification/{token}", h.Verifier) // GET  /api/v1/user/verification/{token}
		r.Get("/{id}", h.Informer)                 // GET  /api/v1/user/{id}

	})
}

func (h *UserHandler) Informer(w http.ResponseWriter, r *http.Request) {

}

func (h *UserHandler) Loginer(w http.ResponseWriter, r *http.Request) {
	tokenResponse := &responses.TokenResponse{}
	loginUser := &dto.LoginUser{}

	_, err := request.ReadAllRequestBody(r, loginUser)
	if err != nil {
		tokenResponse.ErrorUpdate(err, http.StatusBadRequest)
		h.Sender.SendResponse(w, tokenResponse)
		return
	}

	token, err := h.UserService.Login(r.Context(), loginUser)

	// Errors
	if err != nil {

		// Logger
		// +err

		switch {
		// Credentials
		case errors.Is(err, users.ErrInvalidCredentials):
			tokenResponse.ErrorUpdate(users.ErrInvalidCredentials, http.StatusUnauthorized)

		// Verification
		case errors.Is(err, users.ErrVerification):
			tokenResponse.InfoUpdate(vars.MessNeedVerifyAccount, http.StatusForbidden)
		case errors.Is(err, users.ErrAttemptsVerification):
			tokenResponse.ErrorUpdate(users.ErrAttemptsVerification, http.StatusTooManyRequests)

		// Server
		case errors.Is(err, users.ErrServer):
			// users.ErrServer
			tokenResponse.ErrorUpdate(err, http.StatusInternalServerError)

		default:
			tokenResponse.ErrorUpdate(err, http.StatusBadRequest)
		}
		h.Sender.SendResponse(w, tokenResponse)
		return
	}

	tokenResponse.AttrsUpdate(token, http.StatusOK)
	h.Sender.SendResponse(w, tokenResponse)

	_ = token.Refresh

	// Сделать Куки и положить туда!
	// token.Refresh

}

func (h *UserHandler) Verifier(w http.ResponseWriter, r *http.Request) {
	infoResponse := &responses.InfoResponse{}
	token := chi.URLParam(r, Token)

	_, err := h.UserService.Verification(r.Context(), token)

	// Errors
	if err != nil {
		switch {
		// Verification
		case errors.Is(err, users.ErrLinkVerification):
			infoResponse.ErrorUpdate(users.ErrLinkVerification, http.StatusNotFound)
		case errors.Is(err, users.ErrRepeatVerification):
			infoResponse.ErrorUpdate(users.ErrRepeatVerification, http.StatusBadRequest)
		case errors.Is(err, users.ErrAttemptsVerification):
			infoResponse.ErrorUpdate(users.ErrAttemptsVerification, http.StatusTooManyRequests)

		// Server
		case errors.Is(err, users.ErrServer):
			infoResponse.ErrorUpdate(users.ErrServer, http.StatusInternalServerError)

		default:
			infoResponse.ErrorUpdate(err, http.StatusBadRequest)
		}
		h.Sender.SendResponse(w, infoResponse)
		return
	}

	infoResponse.InfoUpdate(vars.MessOkVerification, http.StatusOK)
	h.Sender.SendResponse(w, infoResponse)
}

// Убрать из сервиса подготовку user Response и перенести ее сюда
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	infoResponse := &responses.InfoResponse{}
	createUser := &dto.CreateUser{}

	_, err := request.ReadAllRequestBody(r, createUser)

	if err != nil {
		infoResponse.ErrorUpdate(err, http.StatusBadRequest)
		h.Sender.SendResponse(w, infoResponse)
		return
	}

	// Service
	userDomen, err := h.UserService.Create(r.Context(), createUser)

	// Errors
	if err != nil {
		switch {
		// Credentials
		case errors.Is(err, users.ErrInvalidCredentials):
			infoResponse.ErrorUpdate(users.ErrInvalidCredentials, http.StatusUnauthorized)

			// Server
		case errors.Is(err, users.ErrServer):
			infoResponse.ErrorUpdate(users.ErrServer, http.StatusInternalServerError)

		default:
			infoResponse.ErrorUpdate(err, http.StatusBadRequest)
		}
		h.Sender.SendResponse(w, infoResponse)
		return
	}

	msg := vars.FuncNeedRegistration(userDomen.Email, int(models.VerifTokenLifetime.Minutes()))
	infoResponse.InfoUpdate(msg, http.StatusOK)
	h.Sender.SendResponse(w, infoResponse)
}

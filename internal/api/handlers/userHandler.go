package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/boginskiy/psychologistAI/cmd/config"
	"github.com/boginskiy/psychologistAI/internal/adapters/dto"
	"github.com/boginskiy/psychologistAI/internal/api"
	"github.com/boginskiy/psychologistAI/internal/api/vars"
	"github.com/boginskiy/psychologistAI/internal/errs/server"
	"github.com/boginskiy/psychologistAI/internal/errs/users"
	"github.com/boginskiy/psychologistAI/internal/models/responses"
	"github.com/boginskiy/psychologistAI/internal/service"

	"github.com/boginskiy/psychologistAI/pkg/cookie"
	"github.com/boginskiy/psychologistAI/pkg/request"
	"github.com/go-chi/chi"
)

const Token = "token"

type UserHandler struct {
	UserService    service.UserService
	Cooker         cookie.Cooker
	ResponseSender api.ResponseSender
	basepath       string
}

func NewUserHandler(bpath string, userServ service.UserService, resSender api.ResponseSender, cooker cookie.Cooker) *UserHandler {
	return &UserHandler{
		UserService:    userServ,
		Cooker:         cooker,
		ResponseSender: resSender,
		basepath:       bpath,
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
	infoResponse := &responses.InfoResponse{}
	loginUser := &dto.LoginUser{}

	// Read body
	_, err := request.ReadAllRequestBody(r, loginUser)
	if err != nil {
		infoResponse.ErrorUpdate(err, http.StatusBadRequest)
		h.ResponseSender.SendResponse(w, infoResponse)
		return
	}

	// Take context (Binding). Берем Реальный IP пользователя.
	loginUser.IP = request.TakeRealUserIP(r)
	// Берем инфо с User-Agent. Info: OS, Browser, Device
	loginUser.UserAgent = request.TakeInfoAboutUserAgent(r)

	// Service
	token, err := h.UserService.Login(r.Context(), loginUser)

	// Errors
	if err != nil {

		// Logger
		// +err

		switch {
		// Credentials
		case errors.Is(err, users.ErrInvalidCredentials):
			infoResponse.ErrorUpdate(users.ErrInvalidCredentials, http.StatusUnauthorized)

		// Verification
		case errors.Is(err, users.ErrVerification):
			infoResponse.InfoUpdate(vars.MessNeedVerifyAccount, http.StatusForbidden)
		case errors.Is(err, users.ErrAttemptsVerification):
			infoResponse.ErrorUpdate(users.ErrAttemptsVerification, http.StatusTooManyRequests)

		// Server
		case errors.Is(err, server.ErrServer):
			// users.ErrServer
			infoResponse.ErrorUpdate(err, http.StatusInternalServerError)

		default:
			infoResponse.ErrorUpdate(err, http.StatusBadRequest)
		}
		h.ResponseSender.SendResponse(w, infoResponse)
		return
	}

	// Cookies
	cookieAccessToken, err1 := h.Cooker.CreateCookie(config.COOKIE_NAME_ACCESS_TOKEN, token.Access)
	cookieRefreshToken, err2 := h.Cooker.CreateCookie(config.COOKIE_NAME_REFRESH_TOKEN, token.Refresh)

	if err1 != nil || err2 != nil {
		// + Logger full error
		fmt.Println(fmt.Errorf("%s:%s:%s", server.ErrServer, err1, err2))
		infoResponse.ErrorUpdate(server.ErrServer, http.StatusInternalServerError)
		h.ResponseSender.SendResponse(w, infoResponse)
		return
	}

	// Response
	h.ResponseSender.AddSetCookies(w, cookieAccessToken, cookieRefreshToken)
	infoResponse.InfoUpdate(vars.MessOkLogin, http.StatusOK)
	h.ResponseSender.SendResponse(w, infoResponse)
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
		case errors.Is(err, server.ErrServer):
			infoResponse.ErrorUpdate(server.ErrServer, http.StatusInternalServerError)

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
	infoResponse := &responses.InfoResponse{}
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
		case errors.Is(err, users.ErrInvalidCredentials):
			infoResponse.ErrorUpdate(users.ErrInvalidCredentials, http.StatusUnauthorized)

			// Server
		case errors.Is(err, server.ErrServer):
			infoResponse.ErrorUpdate(server.ErrServer, http.StatusInternalServerError)

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

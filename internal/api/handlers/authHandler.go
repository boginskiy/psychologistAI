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
	"github.com/boginskiy/psychologistAI/internal/models/response"
	"github.com/boginskiy/psychologistAI/internal/service"
	"github.com/boginskiy/psychologistAI/pkg/cookie"
	"github.com/boginskiy/psychologistAI/pkg/request"
	"github.com/go-chi/chi"
)

type AuthHandler struct {
	AuthService    service.AuthService
	Cooker         cookie.Cooker
	ResponseSender api.ResponseSender
	basepath       string
}

func NewAuthHandler(bpath string, authServ service.AuthService, resSender api.ResponseSender, cooker cookie.Cooker) *AuthHandler {
	return &AuthHandler{
		AuthService:    authServ,
		Cooker:         cooker,
		ResponseSender: resSender,
		basepath:       bpath,
	}
}

func (h *AuthHandler) Registration(r chi.Router) {
	r.Route(h.basepath, func(r chi.Router) {
		r.Post("/login", h.Loginer)     // POST /auth/login
		r.Post("/refresh", h.Refresher) // POST /auth/refresh
		r.Post("/logout", h.Logouter)   // POST /auth/logout
	})
}

func (h *AuthHandler) Logouter(w http.ResponseWriter, r *http.Request) {

}

func (h *AuthHandler) Refresher(w http.ResponseWriter, r *http.Request) {
	infoResponse := &response.InfoResponse{}
	refreshTokenRequest := &dto.RefreshTokenRequest{}

	// Берем cookie с refresh token
	cookie, err := r.Cookie(config.COOKIE_NAME_REFRESH_TOKEN)
	if err != nil {
		fmt.Println(fmt.Errorf("%s:%s", users.ErrAuth, err))
		infoResponse.ErrorUpdate(users.ErrAuth, http.StatusUnauthorized)
		h.ResponseSender.SendResponse(w, infoResponse)
	}

	refreshTokenRequest.UserAgent = request.TakeInfoAboutUserAgent(r)
	refreshTokenRequest.IP = request.TakeRealUserIP(r)
	refreshTokenRequest.Token = cookie.Value

	newToken, err := h.AuthService.Refresh(r.Context(), refreshTokenRequest)

}

func (h *AuthHandler) Loginer(w http.ResponseWriter, r *http.Request) {
	infoResponse := &response.InfoResponse{}
	loginUser := &dto.LoginUser{}

	// Read body
	_, err := request.ReadAllRequestBody(r, loginUser)
	if err != nil {
		infoResponse.ErrorUpdate(err, http.StatusBadRequest)
		h.ResponseSender.SendResponse(w, infoResponse)
		return
	}

	// // Take context (Binding). Берем Реальный IP пользователя.
	// loginUser.IP = request.TakeRealUserIP(r)
	// // Берем инфо с User-Agent. Info: OS, Browser, Device
	// loginUser.UserAgent = request.TakeInfoAboutUserAgent(r)

	// Service
	token, err := h.AuthService.Login(r.Context(), loginUser)

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

// TODO:
// Проверка в Мидлвари JWT токена
// Разные сценарии
// Как отправляется куки с access токеном и как будет отправляться с refresh

// Logout
// 3. Ротация и отзыв (Revocation) Так как JWT сам по себе живет своей жизнью до истечения срока,
// тебе обязательно нужен механизм принудительного завершения сессии (например, кнопка
// «Выйти на всех устройствах»).

// При создании refresh token генерируй уникальный идентификатор — JTI (JWT ID).
// Храни этот JTI в быстрой базе (Redis) вместе с UserID и временем жизни.
// В самом payload refresh token тоже положи этот JTI.
// При каждом вызове /auth/refresh проверяй наличие этого JTI в Redis. Если его нет (пользователь нажал Logout ранее) — отклоняй запрос.
// При нажатии /auth/logout удаляй текущий JTI из Redis и присылай Set-Cookie с Max-Age: -1 для удаления обеих кук.

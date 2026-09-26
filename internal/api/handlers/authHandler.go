package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/boginskiy/psychologistAI/cmd/config"
	"github.com/boginskiy/psychologistAI/internal/adapters/dto"
	"github.com/boginskiy/psychologistAI/internal/api"
	"github.com/boginskiy/psychologistAI/internal/api/request"
	"github.com/boginskiy/psychologistAI/internal/api/vars"
	"github.com/boginskiy/psychologistAI/internal/errs"
	"github.com/boginskiy/psychologistAI/internal/middleware"

	"github.com/boginskiy/psychologistAI/internal/models/response"
	"github.com/boginskiy/psychologistAI/internal/service"
	"github.com/boginskiy/psychologistAI/pkg/cookie"
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

func (h *AuthHandler) Registration(r chi.Router, middleware middleware.HandleMiddleware) {
	r.Route(h.basepath, func(r chi.Router) {
		// Public
		r.Group(func(r chi.Router) {
			r.Post("/login", h.Loginer)     // POST /auth/login
			r.Post("/refresh", h.Refresher) // POST /auth/refresh
		})

		// Need Auth
		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(h.AuthService))
			r.Post("/logout", h.Logouter) // POST /auth/logout
		})

	})
}

func (h *AuthHandler) Logouter(w http.ResponseWriter, r *http.Request) {

}

func (h *AuthHandler) Refresher(w http.ResponseWriter, r *http.Request) {
	infoResponse := &response.InfoResponse{}
	refreshTokenRequest := &dto.RefreshTokenRequest{}

	// Take cookie with refresh token
	cookie, err := r.Cookie(config.COOKIE_NAME_REFRESH_TOKEN)
	if err != nil {
		fmt.Println(fmt.Errorf("%s:%s", errs.ErrAuth, err))
		infoResponse.ErrorUpdate(errs.ErrAuth, http.StatusUnauthorized)
		h.ResponseSender.SendResponse(w, infoResponse)
	}

	// Info current request
	refreshTokenRequest.OS, refreshTokenRequest.Browser, refreshTokenRequest.Device = request.TakeDeviceInfo(r)
	refreshTokenRequest.UserAgent = request.TakeUserAgent(r)
	refreshTokenRequest.IP = request.TakeRealUserIP(r)
	refreshTokenRequest.Token = cookie.Value

	// Service
	newToken, err := h.AuthService.Refresh(r.Context(), refreshTokenRequest)

	if err != nil {
		switch {
		case errors.Is(err, errs.ErrSession):
			// + logger
			infoResponse.ErrorUpdate(errs.ErrAuth, http.StatusUnauthorized)

		case errors.Is(err, errs.ErrUsingToken), errors.Is(err, errs.ErrLegitimacyUser),
			errors.Is(err, errs.ErrTimeLiveSession), errors.Is(err, errs.ErrCompareToken):
			// + logger
			infoResponse.ErrorUpdate(errs.ErrAuth, http.StatusUnauthorized)

		case errors.Is(err, errs.ErrServer), errors.Is(err, errs.ErrUpdateDBAfterRefresh):
			// + logger
			infoResponse.ErrorUpdate(errs.ErrServer, http.StatusInternalServerError)

		default:
			infoResponse.ErrorUpdate(err, http.StatusBadRequest)
		}
		h.ResponseSender.SendResponse(w, infoResponse)
		return
	}

	// Cookies

	// Обнуляем текущий  cookie
	oldCookie, err := h.Cooker.ClearCookie(cookie)
	if err != nil {
		// + logger
	}

	cookieAccessToken, err1 := h.Cooker.CreateCookie(config.COOKIE_NAME_ACCESS_TOKEN, newToken.AccessToken)
	cookieRefreshToken, err2 := h.Cooker.CreateCookie(config.COOKIE_NAME_REFRESH_TOKEN, newToken.RefreshToken)

	if err1 != nil || err2 != nil {
		// + Logger
		fmt.Println(fmt.Errorf("%s:%s:%s", errs.ErrServer, err1, err2))
		infoResponse.ErrorUpdate(errs.ErrServer, http.StatusInternalServerError)
		h.ResponseSender.SendResponse(w, infoResponse)
		return
	}

	// Response
	h.ResponseSender.AddSetCookies(w, oldCookie, cookieAccessToken, cookieRefreshToken)
	infoResponse.InfoUpdate(vars.MessOkLogin, http.StatusOK)
	h.ResponseSender.SendResponse(w, infoResponse)
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

	// Info current request
	loginUser.OS, loginUser.Browser, loginUser.Device = request.TakeDeviceInfo(r)
	loginUser.IP = request.TakeRealUserIP(r)
	loginUser.UserAgent = request.TakeUserAgent(r)

	// Service
	token, err := h.AuthService.Login(r.Context(), loginUser)

	// Errors
	if err != nil {

		// Logger
		// +err

		switch {
		// Credentials
		case errors.Is(err, errs.ErrInvalidCredentials):
			infoResponse.ErrorUpdate(errs.ErrInvalidCredentials, http.StatusUnauthorized)

		// Verification
		case errors.Is(err, errs.ErrVerification):
			infoResponse.InfoUpdate(vars.MessNeedVerifyAccount, http.StatusForbidden)
		case errors.Is(err, errs.ErrAttemptsVerification):
			infoResponse.ErrorUpdate(errs.ErrAttemptsVerification, http.StatusTooManyRequests)

		// Server
		case errors.Is(err, errs.ErrServer):
			// errs.ErrServer
			infoResponse.ErrorUpdate(err, http.StatusInternalServerError)

		default:
			infoResponse.ErrorUpdate(err, http.StatusBadRequest)
		}
		h.ResponseSender.SendResponse(w, infoResponse)
		return
	}

	// Cookies
	cookieAccessToken, err1 := h.Cooker.CreateCookie(config.COOKIE_NAME_ACCESS_TOKEN, token.AccessToken)
	cookieRefreshToken, err2 := h.Cooker.CreateCookie(config.COOKIE_NAME_REFRESH_TOKEN, token.RefreshToken)

	if err1 != nil || err2 != nil {
		// + Logger full error
		fmt.Println(fmt.Errorf("%s:%s:%s", errs.ErrServer, err1, err2))
		infoResponse.ErrorUpdate(errs.ErrServer, http.StatusInternalServerError)
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

// Logout
// 3. Ротация и отзыв (Revocation) Так как JWT сам по себе живет своей жизнью до истечения срока,
// тебе обязательно нужен механизм принудительного завершения сессии (например, кнопка
// «Выйти на всех устройствах»).

// При создании refresh token генерируй уникальный идентификатор — JTI (JWT ID).
// Храни этот JTI в быстрой базе (Redis) вместе с UserID и временем жизни.
// В самом payload refresh token тоже положи этот JTI.
// При каждом вызове /auth/refresh проверяй наличие этого JTI в Redis. Если его нет (пользователь нажал Logout ранее) — отклоняй запрос.
// При нажатии /auth/logout удаляй текущий JTI из Redis и присылай Set-Cookie с Max-Age: -1 для удаления обеих кук.

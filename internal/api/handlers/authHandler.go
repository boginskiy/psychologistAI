package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/boginskiy/psychologistAI/cmd/config"
	"github.com/boginskiy/psychologistAI/internal/adapters/dto"
	"github.com/boginskiy/psychologistAI/internal/api"

	"github.com/boginskiy/psychologistAI/internal/api/adapters"
	"github.com/boginskiy/psychologistAI/internal/api/vars"
	"github.com/boginskiy/psychologistAI/internal/errs"
	"github.com/boginskiy/psychologistAI/internal/middleware"

	"github.com/boginskiy/psychologistAI/internal/api/response"
	"github.com/boginskiy/psychologistAI/internal/service"
	"github.com/boginskiy/psychologistAI/pkg/cookie"
	"github.com/go-chi/chi"
)

type AuthHandler struct {
	AuthService service.AuthService
	Cooker      cookie.Cooker
	Responder   api.Responder
	Requester   api.Requester
	basepath    string
}

func NewAuthHandler(
	bpath string,
	authServ service.AuthService,
	responder api.Responder,
	requester api.Requester,
	cooker cookie.Cooker,
) *AuthHandler {

	return &AuthHandler{
		AuthService: authServ,
		Cooker:      cooker,
		Responder:   responder,
		Requester:   requester,
		basepath:    bpath,
	}
}

func (h *AuthHandler) Registration(r chi.Router, middleware middleware.HandleMiddleware) {
	r.Route(h.basepath, func(r chi.Router) {
		// Public
		r.Group(func(r chi.Router) {
			r.Post("/login", h.Loginer) // POST /auth/login
			// r.Post("/refresh", h.Refresher) // POST /auth/refresh
		})

		// Need Auth
		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthApiMiddleware(h.AuthService))
			r.Get("/logout", h.Logouter) // POST /auth/logout
		})
	})
}

func (h *AuthHandler) Logouter(w http.ResponseWriter, r *http.Request) {
	// Service
	err := h.AuthService.Logout(r.Context())
	if err != nil {
		body := response.NewInfoBodyWithErr(errs.ErrInfoContext, http.StatusInternalServerError)
		h.Responder.SendResponse(w, body)
		return
	}

	// Cookies
	oldCookie, err := h.Cooker.ClearCookie(config.COOKIE_NAME_ACCESS_TOKEN)
	if err != nil {
		// + logger
		fmt.Println(fmt.Errorf("%v", err))
	}

	// Response
	h.Responder.AddSetCookies(w, oldCookie)
	body := response.NewInfoBody(vars.MessOkLogout, http.StatusOK)
	h.Responder.SendResponse(w, body)
}

func (h *AuthHandler) Refresher(w http.ResponseWriter, r *http.Request) {
	// Take cookie with refresh token
	cookie, err := r.Cookie(config.COOKIE_NAME_REFRESH_TOKEN)
	if err != nil {
		// + logger
		fmt.Println(fmt.Errorf("%s:%s", errs.ErrAuth, err))

		body := response.NewInfoBodyWithErr(errs.ErrAuth, http.StatusUnauthorized)
		h.Responder.SendResponse(w, body)
		return
	}

	// Put Info from current request
	refreshTokenReq := adapters.NewTokenReq(h.Requester, r, cookie)
	// Service
	newToken, err := h.AuthService.Refresh(r.Context(), refreshTokenReq)

	if err != nil {
		body := &response.InfoBody{}
		// + logger
		fmt.Println(fmt.Errorf("%v", err))

		switch {
		// Credentials
		case errors.Is(err, errs.ErrSession):
			body.ErrorUpdate(errs.ErrAuth, http.StatusUnauthorized)

		case errors.Is(err, errs.ErrUsingToken), errors.Is(err, errs.ErrLegitimacyUser):
			body.ErrorUpdate(errs.ErrAuth, http.StatusUnauthorized)

		case errors.Is(err, errs.ErrTimeLiveSession), errors.Is(err, errs.ErrCompareToken):
			body.ErrorUpdate(errs.ErrAuth, http.StatusUnauthorized)

		// Server
		case errors.Is(err, errs.ErrServer), errors.Is(err, errs.ErrUpdateDBAfterRefresh):
			body.ErrorUpdate(errs.ErrServer, http.StatusInternalServerError)

		default:
			body = response.NewInfoBodyWithErr(err, http.StatusBadRequest)
		}
		h.Responder.SendResponse(w, body)
		return
	}

	// Cookies

	// Обнуляем текущий  cookie
	oldCookie, err := h.Cooker.ClearCookie(config.COOKIE_NAME_REFRESH_TOKEN)
	if err != nil {
		// + logger
		fmt.Println(fmt.Errorf("%v", err))
	}

	cookieAccessToken, err1 := h.Cooker.CreateCookie(config.COOKIE_NAME_ACCESS_TOKEN, newToken.AccessToken)
	cookieRefreshToken, err2 := h.Cooker.CreateCookie(config.COOKIE_NAME_REFRESH_TOKEN, newToken.RefreshToken)

	if err1 != nil || err2 != nil {
		// + Logger
		fmt.Println(fmt.Errorf("%s:%s:%s", errs.ErrServer, err1, err2))

		body := response.NewInfoBodyWithErr(errs.ErrServer, http.StatusInternalServerError)
		h.Responder.SendResponse(w, body)
		return
	}

	// Response
	h.Responder.AddSetCookies(w, oldCookie, cookieAccessToken, cookieRefreshToken)
	body := response.NewInfoBody(vars.MessOkLogin, http.StatusOK)
	h.Responder.SendResponse(w, body)
}

func (h *AuthHandler) Loginer(w http.ResponseWriter, r *http.Request) {
	loginUser := &dto.LoginUser{}

	// Read body
	_, err := h.Requester.ReadAllRequestBody(r, loginUser)
	if err != nil {
		body := response.NewInfoBodyWithErr(err, http.StatusBadRequest)
		h.Responder.SendResponse(w, body)
		return
	}

	// Put Info from current request
	loginUser = adapters.UpdateLoginUserFromRequest(loginUser, h.Requester, r)
	// Service
	token, err := h.AuthService.Login(r.Context(), loginUser)

	// Errors
	if err != nil {
		body := &response.InfoBody{}
		// + logger
		fmt.Println(fmt.Errorf("%v", err))

		switch {
		// Credentials
		case errors.Is(err, errs.ErrInvalidCredentials):
			body.ErrorUpdate(errs.ErrInvalidCredentials, http.StatusUnauthorized)

		// Verification
		case errors.Is(err, errs.ErrVerification):
			body.InfoUpdate(vars.MessNeedVerifyAccount, http.StatusForbidden)
		case errors.Is(err, errs.ErrAttemptsVerification):
			body.ErrorUpdate(errs.ErrAttemptsVerification, http.StatusTooManyRequests)

		// Server
		case errors.Is(err, errs.ErrServer):
			body.ErrorUpdate(err, http.StatusInternalServerError)

		default:
			body.ErrorUpdate(err, http.StatusBadRequest)
		}
		h.Responder.SendResponse(w, body)
		return
	}

	// Cookies
	cookieAccessToken, err1 := h.Cooker.CreateCookie(config.COOKIE_NAME_ACCESS_TOKEN, token.AccessToken)
	cookieRefreshToken, err2 := h.Cooker.CreateCookie(config.COOKIE_NAME_REFRESH_TOKEN, token.RefreshToken)

	if err1 != nil || err2 != nil {
		// + Logger
		fmt.Println(fmt.Errorf("%s:%s:%s", errs.ErrServer, err1, err2))

		body := response.NewInfoBodyWithErr(errs.ErrServer, http.StatusInternalServerError)
		h.Responder.SendResponse(w, body)
		return
	}

	// Response
	h.Responder.AddSetCookies(w, cookieAccessToken, cookieRefreshToken)
	body := response.NewInfoBody(vars.MessOkLogin, http.StatusOK)
	h.Responder.SendResponse(w, body)
}

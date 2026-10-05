package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/boginskiy/psychologistAI/cmd/config"
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

type RefreshHandler struct {
	AuthService service.AuthService
	Cooker      cookie.Cooker
	Responder   api.Responder
	Requester   api.Requester
}

func NewRefreshHandler(

	authServ service.AuthService,
	responder api.Responder,
	requester api.Requester,
	cooker cookie.Cooker,
) *RefreshHandler {

	return &RefreshHandler{
		AuthService: authServ,
		Cooker:      cooker,
		Responder:   responder,
		Requester:   requester,
	}
}

func (h *RefreshHandler) Registration(r chi.Router, middleware middleware.HandleMiddleware) {
	r.Use(middleware.AuthApiMiddleware(h.AuthService))
	r.Post("/refresh", h.Refresher)
}

func (h *RefreshHandler) Refresher(w http.ResponseWriter, r *http.Request) {
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

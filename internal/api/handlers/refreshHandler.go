package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/boginskiy/psychologistAI/cmd/config"
	"github.com/boginskiy/psychologistAI/internal/adapters"
	"github.com/boginskiy/psychologistAI/internal/api"

	adaptersAPI "github.com/boginskiy/psychologistAI/internal/api/adapters"
	"github.com/boginskiy/psychologistAI/internal/errs"
	"github.com/boginskiy/psychologistAI/internal/middleware"
	"github.com/boginskiy/psychologistAI/internal/msgs"

	"github.com/boginskiy/psychologistAI/internal/api/response"
	"github.com/boginskiy/psychologistAI/internal/service"
	"github.com/boginskiy/psychologistAI/pkg/cooker"
	"github.com/go-chi/chi"
)

type RefreshHandler struct {
	AuthService service.AuthService
	Cooker      cooker.Cooker
	Responder   api.Responder
	Requester   api.Requester
}

func NewRefreshHandler(

	authServ service.AuthService,
	responder api.Responder,
	requester api.Requester,
	cooker cooker.Cooker,
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
		log.Printf("error: %s:%s\n", errs.ErrAuth, err) // + logger

		body := response.NewInfoBodyWithErr(errs.ErrAuth, http.StatusUnauthorized)
		h.Responder.SendResponse(w, body)
		return
	}

	// Put Info from current request
	refreshTokenReq := adaptersAPI.NewTokenReq(h.Requester, r, cookie)
	// Service
	tokenPair, err := h.AuthService.Refresh(r.Context(), refreshTokenReq)

	if err != nil {
		log.Printf("error: %v\n", err) // + logger

		body := &response.InfoBody{}

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

	// Zeroing out cookie
	oldCookie, err := h.Cooker.ClearCookie(config.COOKIE_NAME_REFRESH_TOKEN)
	if err != nil {
		log.Printf("error: %v\n", err) // + logger
	}

	// Create cookie with tokens
	cookiePair, err := adapters.FromTokenPair(h.Cooker, tokenPair)
	if err != nil {
		log.Printf("error: %v\n", err) // + logger

		body := response.NewInfoBodyWithErr(errs.ErrServer, http.StatusInternalServerError)
		h.Responder.SendResponse(w, body)
		return
	}

	// Response
	h.Responder.AddSetCookies(w, oldCookie, cookiePair.Access, cookiePair.Refresh)
	body := response.NewInfoBody(msgs.MessOkLogin, http.StatusOK)
	h.Responder.SendResponse(w, body)
}

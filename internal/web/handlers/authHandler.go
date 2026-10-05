package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/boginskiy/psychologistAI/cmd/config"
	"github.com/boginskiy/psychologistAI/internal/api"
	adaptersApi "github.com/boginskiy/psychologistAI/internal/api/adapters"
	"github.com/boginskiy/psychologistAI/internal/errs"
	"github.com/boginskiy/psychologistAI/internal/middleware"
	"github.com/boginskiy/psychologistAI/internal/service"
	"github.com/boginskiy/psychologistAI/internal/web/adapters"
	adaptersWeb "github.com/boginskiy/psychologistAI/internal/web/adapters"
	"github.com/boginskiy/psychologistAI/internal/web/renders"
	"github.com/boginskiy/psychologistAI/pkg/cookie"
	"github.com/go-chi/chi"
)

type AuthHandler struct {
	Requester   api.Requester
	Responder   api.Responder
	Cooker      cookie.Cooker
	AuthService service.AuthService
}

func NewAuthHandler(requester api.Requester, responder api.Responder, authService service.AuthService, cookie cookie.Cooker) *AuthHandler {
	return &AuthHandler{
		Requester:   requester,
		Responder:   responder,
		AuthService: authService,
		Cooker:      cookie,
	}
}

func (h *AuthHandler) Registration(r chi.Router, middleware middleware.HandleMiddleware) {
	r.Get("/login", h.ShowLoginForm)
	r.Post("/login", h.Loginer)
}

func (h *AuthHandler) ShowLoginForm(w http.ResponseWriter, r *http.Request) {
	// Рендерим шаблон login
	tmpl, err := template.ParseFiles(
		"templates/base.html",
		"templates/login.html",
	)
	if err != nil {
		renders.RenderError(w, "500", http.StatusInternalServerError, nil)
		return
	}

	tmpl.ExecuteTemplate(w, "base", nil)
}

func (h *AuthHandler) Loginer(w http.ResponseWriter, r *http.Request) {
	loginUser, err := adaptersWeb.ToLoginUserFromFormRequest(r)

	if err != nil {
		// + logger
		log.Printf("template error: %v", err)
		renders.RenderError(w, "400", http.StatusBadRequest, nil)
		return
	}

	// Put Info from current request
	loginUser = adaptersApi.UpdateLoginUserFromRequest(loginUser, h.Requester, r)
	// Service
	token, err := h.AuthService.Login(r.Context(), loginUser)

	// Errors
	if err != nil {
		// + logger
		fmt.Println(fmt.Errorf("%v", err))

		switch {
		// Credentials
		case errors.Is(err, errs.ErrInvalidCredentials):
			renders.RenderError(w, "401", http.StatusUnauthorized, nil)

		// Verification
		case errors.Is(err, errs.ErrVerification):
			renders.RenderError(w, "403", http.StatusForbidden, nil)
		case errors.Is(err, errs.ErrAttemptsVerification):
			renders.RenderError(w, "429", http.StatusTooManyRequests, nil)

		// Server
		case errors.Is(err, errs.ErrServer):
			renders.RenderError(w, "500", http.StatusInternalServerError, nil)

		default:
			renders.RenderError(w, "400", http.StatusBadRequest, nil)
		}

		return
	}

	// Cookies
	cookieAccessToken, err1 := h.Cooker.CreateCookie(config.COOKIE_NAME_ACCESS_TOKEN, token.AccessToken)
	cookieRefreshToken, err2 := h.Cooker.CreateCookie(config.COOKIE_NAME_REFRESH_TOKEN, token.RefreshToken)

	if err1 != nil || err2 != nil {
		// + Logger
		fmt.Println(fmt.Errorf("%s:%s:%s", errs.ErrServer, err1, err2))

		renders.RenderError(w, "500", http.StatusInternalServerError, nil)
		return
	}

	// Response
	h.Responder.AddSetCookies(w, cookieAccessToken, cookieRefreshToken)

	startTemplate := adapters.ToMapStartTemplate(true)

	tmpl, err := template.ParseFiles(
		"templates/base.html",
		"templates/index.html",
		"templates/chat.html",
	)

	if err != nil {
		// + logger
		log.Printf("template error: %v", err)

		renders.RenderError(w, "500", http.StatusInternalServerError, nil)
		return
	}

	tmpl.ExecuteTemplate(w, "base", startTemplate)
}

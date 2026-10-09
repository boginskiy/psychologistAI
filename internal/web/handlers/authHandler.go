package handlers

import (
	"errors"
	"log"
	"net/http"
	"text/template"

	"github.com/boginskiy/psychologistAI/internal/adapters"
	"github.com/boginskiy/psychologistAI/internal/api"
	adaptersApi "github.com/boginskiy/psychologistAI/internal/api/adapters"
	"github.com/boginskiy/psychologistAI/internal/errs"
	"github.com/boginskiy/psychologistAI/internal/middleware"
	"github.com/boginskiy/psychologistAI/internal/service"
	adaptersWeb "github.com/boginskiy/psychologistAI/internal/web/adapters"
	"github.com/boginskiy/psychologistAI/internal/web/renders"
	"github.com/boginskiy/psychologistAI/pkg/cooker"
	"github.com/go-chi/chi"
)

type AuthHandler struct {
	Requester   api.Requester
	Responder   api.Responder
	Cooker      cooker.Cooker
	AuthService service.AuthService
}

func NewAuthHandler(requester api.Requester, responder api.Responder, authService service.AuthService, cooker cooker.Cooker) *AuthHandler {
	return &AuthHandler{
		Requester:   requester,
		Responder:   responder,
		AuthService: authService,
		Cooker:      cooker,
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
	tokenPair, err := h.AuthService.Login(r.Context(), loginUser)

	// Errors
	if err != nil {
		log.Printf("error: %v\n", err) // + logger

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

	// Create cookie with tokens
	cookiePair, err := adapters.FromTokenPair(h.Cooker, tokenPair)
	if err != nil {
		log.Printf("error: %v\n", err) // + logger

		renders.RenderError(w, "500", http.StatusInternalServerError, nil)
		return
	}

	// Response
	h.Responder.AddSetCookies(w, cookiePair.Access, cookiePair.Refresh)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

package handlers

import (
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
	Catcher     errs.Catcher
}

func NewAuthHandler(requester api.Requester, responder api.Responder, authService service.AuthService, cooker cooker.Cooker) *AuthHandler {
	return &AuthHandler{
		Requester:   requester,
		Responder:   responder,
		AuthService: authService,
		Cooker:      cooker,
		Catcher:     errs.NewErrService(responder),
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
	// Catch error
	if err != nil {
		h.Catcher.CatchWebErr(w, err)
		return
	}

	// Cookies. Create cookie with tokens
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

package handlers

import (
	"log"
	"net/http"
	"text/template"

	"github.com/boginskiy/psychologistAI/internal/api/request"
	"github.com/boginskiy/psychologistAI/internal/middleware"
	"github.com/boginskiy/psychologistAI/internal/service"
	"github.com/boginskiy/psychologistAI/internal/web/adapters"
	"github.com/boginskiy/psychologistAI/internal/web/renders"
	"github.com/go-chi/chi"
)

type HomeHandler struct {
	AuthService service.AuthService
}

func NewHomeHandler(authService service.AuthService) *HomeHandler {
	return &HomeHandler{
		AuthService: authService,
	}
}

func (h *HomeHandler) Registration(r chi.Router, middleware middleware.HandleMiddleware) {
	r.Use(middleware.AuthWebMiddleware(h.AuthService))
	r.Get("/", h.Start)
}

func (h *HomeHandler) Start(w http.ResponseWriter, r *http.Request) {
	_, isUser := request.GetInfoUserFromContext(r.Context())
	startTemplate := adapters.ToMapStartTemplate(isUser)

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

package handlers

import (
	"log"
	"net/http"
	"text/template"

	"github.com/boginskiy/psychologistAI/internal/middleware"
	"github.com/boginskiy/psychologistAI/internal/web/renders"
	"github.com/go-chi/chi"
)

type HomeHandler struct {
	basepath string
}

func NewHomeHandler(
	bpath string,
) *HomeHandler {

	return &HomeHandler{
		basepath: bpath,
	}
}

func (h *HomeHandler) Registration(r chi.Router, middleware middleware.HandleMiddleware) {
	r.Route(h.basepath, func(r chi.Router) {
		r.Get("/", h.Start)
	})
}

func (h *HomeHandler) Start(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"templates/base.html",
		"templates/chat.html",
	)

	if err != nil {
		// + logger
		log.Printf("template error: %v", err)

		renders.RenderError(w, "500", http.StatusInternalServerError, nil)
		return
	}

	tmpl.ExecuteTemplate(w, "base", nil)

}

package router

import (
	"context"
	"net/http"

	"github.com/boginskiy/psychologistAI/internal/api/handlers"
	"github.com/boginskiy/psychologistAI/internal/api/response"
	"github.com/boginskiy/psychologistAI/internal/errs"
	"github.com/boginskiy/psychologistAI/internal/middleware"
	"github.com/boginskiy/psychologistAI/internal/web/renders"
	"github.com/go-chi/chi"
)

type RouterChi struct {
	Mux        *chi.Mux
	Middleware middleware.Middleware
}

func NewRouterChi(ctx context.Context, middlew middleware.Middleware) *RouterChi {
	return &RouterChi{
		Mux:        chi.NewRouter(),
		Middleware: middlew,
	}
}

func (c *RouterChi) Run() http.Handler {
	return c.Mux
}

func (c *RouterChi) RegisterAPIRoutes(start string, handlers ...handlers.Registrar) http.Handler {
	apiMux := chi.NewRouter()

	apiMux.Use(c.Middleware.RecoveryMiddleware)
	apiMux.Use(c.Middleware.LoggingMiddleware)

	// API
	for _, handler := range handlers {
		handler.Registration(apiMux, c.Middleware)
	}

	// Not Found
	apiMux.NotFound(func(w http.ResponseWriter, r *http.Request) {
		body := response.NewInfoBodyWithErr(errs.ErrRequest, http.StatusBadRequest)
		Responder := response.NewResponse()
		Responder.SendResponse(w, body)
	})

	c.Mux.Mount(start, apiMux)
	return c.Mux
}

func (c *RouterChi) RegisterWEBRoutes(start string, handlers ...handlers.Registrar) http.Handler {
	webMux := chi.NewRouter()

	webMux.Use(c.Middleware.RecoveryMiddleware)
	webMux.Use(c.Middleware.LoggingMiddleware)

	// WEB
	for _, handler := range handlers {
		handler.Registration(webMux, c.Middleware)
	}

	// Not Found
	webMux.NotFound(func(w http.ResponseWriter, r *http.Request) {
		renders.RenderError(w, "404", http.StatusNotFound, nil)
	})

	c.Mux.Mount(start, webMux)

	// Раздача статики
	c.Mux.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	return c.Mux
}

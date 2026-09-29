package router

import (
	"context"
	"net/http"

	"github.com/boginskiy/psychologistAI/internal/api/handlers"
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

	apiMux.Route(start, func(r chi.Router) {
		r.Use(c.Middleware.RecoveryMiddleware)
		r.Use(c.Middleware.LoggingMiddleware)

		// API
		for _, handler := range handlers {
			handler.Registration(r, c.Middleware)
		}
	})

	c.Mux.Mount(start, apiMux)
	return c.Mux
}

func (c *RouterChi) RegisterWEBRoutes(start string, handlers ...handlers.Registrar) http.Handler {
	webMux := chi.NewRouter()

	webMux.Route(start, func(r chi.Router) {
		// r.Use(c.Middleware.RecoveryMiddleware)
		// r.Use(c.Middleware.LoggingMiddleware)

		// WEB
		for _, handler := range handlers {
			handler.Registration(r, c.Middleware)
		}

		// Not Found
		r.NotFound(func(w http.ResponseWriter, r *http.Request) {
			renders.RenderError(w, "404", http.StatusNotFound, nil)
		})
	})

	c.Mux.Mount(start, webMux)
	return c.Mux
}

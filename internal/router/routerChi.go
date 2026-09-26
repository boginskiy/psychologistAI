package router

import (
	"context"
	"net/http"

	"github.com/boginskiy/psychologistAI/internal/api/handlers"
	"github.com/boginskiy/psychologistAI/internal/middleware"
	"github.com/go-chi/chi"
)

type RouterChi struct {
	mux        *chi.Mux
	Middleware middleware.Middleware
}

func NewRouterChi(ctx context.Context, middlew middleware.Middleware) *RouterChi {
	return &RouterChi{
		mux:        chi.NewRouter(),
		Middleware: middlew,
	}
}

func (c *RouterChi) Run() http.Handler {
	return c.mux
}

func (c *RouterChi) RegisterRoutes(handlers ...handlers.Registrar) http.Handler {
	c.mux.Route("", func(r chi.Router) {

		// Добавляем middleware для всей API
		r.Use(c.Middleware.RecoveryMiddleware)
		r.Use(c.Middleware.LoggingMiddleware)

		for _, handler := range handlers {
			handler.Registration(r, c.Middleware)
		}
	})

	return c.mux
}

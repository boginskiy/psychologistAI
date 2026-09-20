package router

import (
	"context"
	"net/http"

	"github.com/boginskiy/psychologistAI/internal/api/handlers"
	"github.com/go-chi/chi"
)

type RouterChi struct {
	mux  *chi.Mux
	home string
}

func NewRouterChi(ctx context.Context, home string) *RouterChi {
	return &RouterChi{
		mux:  chi.NewRouter(),
		home: home,
	}
}

func (c *RouterChi) Run() http.Handler {
	return c.mux
}

func (c *RouterChi) RegisterRoutes(handlers ...handlers.Registrar) http.Handler {
	c.mux.Route(c.home, func(r chi.Router) {

		// Добавляем middleware для всей API
		//r.Use(middleware.Logger)

		for _, handler := range handlers {
			handler.Registration(r)
		}
	})

	return c.mux
}

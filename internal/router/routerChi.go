package router

import (
	"context"

	"github.com/boginskiy/psychologistAI/internal/api/handlers"
	"github.com/go-chi/chi"
)

type RouterChi struct {
	R *chi.Mux
}

func NewRouterChi(ctx context.Context) *RouterChi {
	return &RouterChi{
		R: chi.NewRouter(),
	}
}

func (r *RouterChi) Registration(userHandler, homeHandler handlers.Registrar) {
	r.R.Route("/api/v1", func(r chi.Router) {

		// Добавляем middleware для всей API
		// r.Use(middleware.Logger)

		// Регистрируем маршруты.
		r.Route("/users", func(r chi.Router) { userHandler.Registration(r) })
		r.Route("/home", func(r chi.Router) { homeHandler.Registration(r) })
	})
}

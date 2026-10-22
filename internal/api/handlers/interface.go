package handlers

import (
	"github.com/boginskiy/psychologistAI/internal/middleware"
	"github.com/go-chi/chi"
)

type Registrar interface {
	Registration(r chi.Router, middleware middleware.HandleMiddleware)
}

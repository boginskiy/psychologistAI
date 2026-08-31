package router

import (
	"net/http"

	"github.com/boginskiy/psychologistAI/internal/api/handlers"
)

type Router interface {
	RegisterRoutes(handlers ...handlers.Registrar) http.Handler
	Run() http.Handler
}

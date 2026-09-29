package router

import (
	"net/http"

	"github.com/boginskiy/psychologistAI/internal/api/handlers"
)

type Router interface {
	RegisterAPIRoutes(start string, handlers ...handlers.Registrar) http.Handler
	RegisterWEBRoutes(start string, handlers ...handlers.Registrar) http.Handler
	Run() http.Handler
}

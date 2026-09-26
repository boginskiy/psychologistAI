package middleware

import (
	"net/http"

	"github.com/boginskiy/psychologistAI/internal/service"
)

type InfraMiddleware interface {
	LoggingMiddleware(next http.Handler) http.Handler
	RecoveryMiddleware(next http.Handler) http.Handler
}

type HandleMiddleware interface {
	AuthMiddleware(service service.AuthService) func(http.Handler) http.Handler
}

type Middleware interface {
	InfraMiddleware
	HandleMiddleware
}

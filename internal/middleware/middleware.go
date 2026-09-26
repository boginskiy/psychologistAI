package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/boginskiy/psychologistAI/internal/api"
	"github.com/boginskiy/psychologistAI/internal/errs"
	"github.com/boginskiy/psychologistAI/internal/models/response"
	"github.com/boginskiy/psychologistAI/internal/service"
)

type Middlew struct {

	// AuthService    service.AuthService
	ResponseSender api.ResponseSender
}

func NewMiddlew(resSender api.ResponseSender) *Middlew {
	return &Middlew{
		ResponseSender: resSender,
	}
}

// LoggingMiddleware - вставить сюда свой logger
func (m *Middlew) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		// + logger
		log.Printf("time server response: %s %s %v", r.Method, r.URL.Path, time.Since(start))
	})
}

// RecoveryMiddleware
func (m *Middlew) RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				// + logger
				log.Printf("panic: %v\n", rec)

				infoResponse := &response.InfoResponse{}
				infoResponse.ErrorUpdate(errs.ErrServer, http.StatusInternalServerError)
				m.ResponseSender.SendResponse(w, infoResponse)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (m *Middlew) AuthMiddleware(service service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// service.Login()

			start := time.Now()
			next.ServeHTTP(w, r)
			log.Printf("Time Server Response: %s %s %v", r.Method, r.URL.Path, time.Since(start))
		})
	}
}

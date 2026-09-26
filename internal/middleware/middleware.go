package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/boginskiy/psychologistAI/cmd/config"
	"github.com/boginskiy/psychologistAI/internal/adapters/dto"
	"github.com/boginskiy/psychologistAI/internal/api"
	"github.com/boginskiy/psychologistAI/internal/api/request"
	"github.com/boginskiy/psychologistAI/internal/errs"
	"github.com/boginskiy/psychologistAI/internal/models/response"
	"github.com/boginskiy/psychologistAI/internal/service"
)

type Middlew struct {
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

func (m *Middlew) AuthMiddleware(authService service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Cookie
			cookie, err := r.Cookie(config.COOKIE_NAME_ACCESS_TOKEN)
			if err != nil {
				// + logger
				log.Printf("%v: %v", errs.ErrAuth, err)

				infoResponse := &response.InfoResponse{}
				infoResponse.ErrorUpdate(errs.ErrAuth, http.StatusUnauthorized)
				m.ResponseSender.SendResponse(w, infoResponse)
				return
			}

			// Request
			accTokenReq := &dto.AccessTokenRequest{}

			accTokenReq.OS, accTokenReq.Browser, accTokenReq.Device = request.TakeDeviceInfo(r)
			accTokenReq.UserAgent = request.TakeUserAgent(r)
			accTokenReq.IP = request.TakeRealUserIP(r)
			accTokenReq.Token = cookie.Value

			// Service
			infoUser, err := authService.Auth(r.Context(), accTokenReq)

			if err != nil {
				// + logger
				log.Printf("%v: %v", errs.ErrAuth, err)

				infoResponse := &response.InfoResponse{}
				infoResponse.ErrorUpdate(errs.ErrAuth, http.StatusUnauthorized)
				m.ResponseSender.SendResponse(w, infoResponse)
				return
			}

			// Context
			newCtx := service.SetInfoUser(r.Context(), *infoUser)
			next.ServeHTTP(w, r.WithContext(newCtx))
		})
	}
}

// Рефакторинг AuthMiddleware, еще раз все пройти посмотреть логику.
// Делать далее /logout

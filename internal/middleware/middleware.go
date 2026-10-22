package middleware

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/boginskiy/psychologistAI/cmd/config"
	"github.com/boginskiy/psychologistAI/internal/api"
	"github.com/boginskiy/psychologistAI/internal/api/adapters"
	"github.com/boginskiy/psychologistAI/internal/api/request"
	"github.com/boginskiy/psychologistAI/internal/api/response"
	"github.com/boginskiy/psychologistAI/internal/errs"
	"github.com/boginskiy/psychologistAI/internal/service"
)

type Middlew struct {
	Responder api.Responder
	Requester api.Requester
}

func NewMiddlew(responder api.Responder, requester api.Requester) *Middlew {
	return &Middlew{
		Responder: responder,
		Requester: requester,
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

				body := response.NewInfoBodyWithErr(errs.ErrServer, http.StatusInternalServerError)
				m.Responder.SendResponse(w, body)
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
				fmt.Println(fmt.Errorf("%v:%v", errs.ErrAuth, err))

				body := response.NewInfoBodyWithErr(errs.ErrAuth, http.StatusUnauthorized)
				m.Responder.SendResponse(w, body)
				return
			}

			// Put Info from current request
			accessTokenReq := adapters.NewTokenReq(m.Requester, r, cookie)
			// Service
			infoUser, err := authService.Auth(r.Context(), accessTokenReq)

			if err != nil {
				// + logger
				fmt.Println(fmt.Errorf("%v:%v", errs.ErrAuth, err))

				body := response.NewInfoBodyWithErr(errs.ErrAuth, http.StatusUnauthorized)
				m.Responder.SendResponse(w, body)
				return
			}

			// Context
			newCtx := request.SetInfoUserToContext(r.Context(), *infoUser)
			next.ServeHTTP(w, r.WithContext(newCtx))
		})
	}
}

// Посмотреть, что с рефреш, отдельный метод мидвари?
// Доделать рефреш
// Делать далее /logout
// Далее делаем анкету и продумываем фронт

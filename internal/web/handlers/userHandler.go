package handlers

import (
	"errors"
	"html/template"
	"log"
	"net/http"

	"github.com/boginskiy/psychologistAI/cmd/config"
	"github.com/boginskiy/psychologistAI/internal/errs"
	"github.com/boginskiy/psychologistAI/internal/middleware"
	"github.com/boginskiy/psychologistAI/internal/service"
	"github.com/boginskiy/psychologistAI/internal/web/adapters"
	"github.com/boginskiy/psychologistAI/internal/web/renders"
	"github.com/go-chi/chi"
)

type UserHandler struct {
	basepath    string
	UserService service.UserService
}

func NewUserHandler(
	bpath string,
	userService service.UserService,
) *UserHandler {

	return &UserHandler{
		basepath:    bpath,
		UserService: userService,
	}
}

func (h *UserHandler) Registration(r chi.Router, middleware middleware.HandleMiddleware) {
	r.Route(h.basepath, func(r chi.Router) {
		// Public
		r.Group(func(r chi.Router) {
			r.Get("/registration", h.ShowRegistrForm)
			r.Post("/registration", h.Register)
			r.Get("/verification/{token}", h.Verifier)
		})
	})
}

func (h *UserHandler) Verifier(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")

	// Service
	_, err := h.UserService.Verification(r.Context(), token)

	if err != nil {
		// + logger
		log.Printf("template error: %v", err)
		renders.RenderError(w, "401", http.StatusUnauthorized, nil) //Окно в тем чтобы переотрпвить письмо про авторизацию
		return
	}

	tmpl, err := template.ParseFiles(
		"templates/base.html",
		"templates/verification.html",
	)

	tmpl.ExecuteTemplate(w, "base", nil)
}

func (h *UserHandler) ShowRegistrForm(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"templates/base.html",
		"templates/registration/registr_form.html",
	)
	if err != nil {
		// + logger
		log.Printf("template error: %v", err)
		renders.RenderError(w, "500", http.StatusInternalServerError, nil)
		return
	}

	tmpl.ExecuteTemplate(w, "base", nil)
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	createUser, err := adapters.ToCreateUserFromFormRequest(r)
	if err != nil {
		// + logger
		log.Printf("template error: %v", err)
		renders.RenderError(w, "400", http.StatusBadRequest, nil)
	}

	// Service
	userDomen, err := h.UserService.Create(r.Context(), createUser)

	// Errors
	if err != nil {
		// + logger
		log.Printf("user service error: %v", err)

		switch {
		// Credentials
		case errors.Is(err, errs.ErrInvalidCredentials):
			renders.RenderError(w, "401", http.StatusUnauthorized, nil)

		// Server
		case errors.Is(err, errs.ErrServer):
			renders.RenderError(w, "500", http.StatusInternalServerError, nil)

		default:
			renders.RenderError(w, "400", http.StatusBadRequest, nil)
		}
		return
	}

	registrTemplate := adapters.ToMapRegistrTemplate(userDomen.Email, config.LIVE_TIME_VARIFICATION_TOKEN)

	tmpl, err := template.ParseFiles(
		"templates/base.html",
		"templates/registration/verify_email.html",
	)

	if err != nil {
		// + logger
		log.Printf("template error: %v", err)
		renders.RenderError(w, "500", http.StatusInternalServerError, nil)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", registrTemplate)

	if err != nil {
		// + logger
		log.Printf("template error: %v", err)
		renders.RenderError(w, "500", http.StatusInternalServerError, nil)
		return
	}
}

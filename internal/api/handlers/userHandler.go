package handlers

import (
	"errors"
	"net/http"

	"github.com/boginskiy/psychologistAI/internal/adapters/dto"
	"github.com/boginskiy/psychologistAI/internal/api"
	"github.com/boginskiy/psychologistAI/internal/api/mess"
	"github.com/boginskiy/psychologistAI/internal/api/response"
	"github.com/boginskiy/psychologistAI/internal/service/errs"
	"github.com/boginskiy/psychologistAI/internal/user"
	"github.com/boginskiy/psychologistAI/internal/user/models"
	"github.com/boginskiy/psychologistAI/pkg/request"
	"github.com/go-chi/chi"
)

type UserHandler struct {
	UserService user.UserService
	Sender      api.Sender
	basepath    string
}

func NewUserHandler(bpath string, userServ user.UserService, sender api.Sender) *UserHandler {
	return &UserHandler{
		UserService: userServ,
		Sender:      sender,
		basepath:    bpath,
	}
}

func (h *UserHandler) Registration(r chi.Router) {
	r.Route(h.basepath, func(r chi.Router) {
		r.Post("/registration", h.Register) // POST /api/v1/user/registration
		r.Post("/login", h.Loginer)         // POST /api/v1/user/login

		r.Get("/verification/{token}", h.Verifier) // GET  /api/v1/user/verification/{token}
		r.Get("{id}", h.Informer)                  // GET  /api/v1/user/{id}

	})
}

func (h *UserHandler) Informer(w http.ResponseWriter, r *http.Request) {
	return
}

func (h *UserHandler) Loginer(w http.ResponseWriter, r *http.Request) {
	messResponse := &response.MessRes{}
	loginUser := &dto.LoginUser{}

	_, err := request.ReadAllRequestBody(r, loginUser)
	if err != nil {
		messResponse.UpdateErr(err, http.StatusBadRequest)
		h.Sender.SendResponse(w, messResponse)
		return
	}

	token, err := h.UserService.Login(r.Context(), loginUser)

	// Errors
	if err != nil {

		// Logger
		// +err

		switch {
		// Credentials
		case errors.Is(err, errs.ErrInvalidCredentials):
			messResponse.UpdateErr(errs.ErrInvalidCredentials, http.StatusUnauthorized)

		// Verification
		case errors.Is(err, errs.ErrVerification):
			messResponse.UpdateInfo(mess.MessNeedVerifyAccount, http.StatusForbidden)
		case errors.Is(err, errs.ErrAttemptsVerification):
			messResponse.UpdateErr(errs.ErrAttemptsVerification, http.StatusTooManyRequests)

		// Server
		case errors.Is(err, errs.ErrServer):
			messResponse.UpdateErr(errs.ErrServer, http.StatusInternalServerError)

		default:
			messResponse.UpdateErr(err, http.StatusBadRequest)
		}
		h.Sender.SendResponse(w, messResponse)
		return
	}

	token.Access

	// Сделать Куки и положить туда!
	token.Refresh

	// Прикрепить JWT
}

func (h *UserHandler) Verifier(w http.ResponseWriter, r *http.Request) {
	messResponse := &response.MessRes{}
	var token string

	_, err := request.ReadAllRequestBody(r, &token)
	if err != nil {
		messResponse.UpdateErr(err, http.StatusBadRequest)
		h.Sender.SendResponse(w, messResponse)
		return
	}

	_, err = h.UserService.Verification(r.Context(), token)

	// Errors
	if err != nil {
		switch {
		// Verification
		case errors.Is(err, errs.ErrLinkVerification):
			messResponse.UpdateErr(errs.ErrLinkVerification, http.StatusNotFound)
		case errors.Is(err, errs.ErrRepeatVerification):
			messResponse.UpdateErr(errs.ErrRepeatVerification, http.StatusBadRequest)
		case errors.Is(err, errs.ErrAttemptsVerification):
			messResponse.UpdateErr(errs.ErrAttemptsVerification, http.StatusTooManyRequests)

		// Server
		case errors.Is(err, errs.ErrServer):
			messResponse.UpdateErr(errs.ErrServer, http.StatusInternalServerError)

		default:
			messResponse.UpdateErr(err, http.StatusBadRequest)
		}
		h.Sender.SendResponse(w, messResponse)
		return
	}

	messResponse.UpdateInfo(mess.MessOkVerification, http.StatusOK)
	h.Sender.SendResponse(w, messResponse)
}

// Убрать из сервиса подготовку user Response и перенести ее сюда
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	messResponse := &response.MessRes{}
	createUser := &dto.CreateUser{}

	_, err := request.ReadAllRequestBody(r, createUser)

	if err != nil {
		messResponse.UpdateErr(err, http.StatusBadRequest)
		h.Sender.SendResponse(w, messResponse)
		return
	}

	// Service
	userDomen, err := h.UserService.Create(r.Context(), createUser)

	// Errors
	if err != nil {
		switch {
		// Credentials
		case errors.Is(err, errs.ErrInvalidCredentials):
			messResponse.UpdateErr(errs.ErrInvalidCredentials, http.StatusUnauthorized)

			// Server
		case errors.Is(err, errs.ErrServer):
			messResponse.UpdateErr(errs.ErrServer, http.StatusInternalServerError)

		default:
			messResponse.UpdateErr(err, http.StatusBadRequest)
		}
		h.Sender.SendResponse(w, messResponse)
		return
	}

	msg := mess.FuncNeedRegistration(userDomen.Email, int(models.VerifTokenLifetime.Minutes()))
	messResponse.UpdateInfo(msg, http.StatusOK)
	h.Sender.SendResponse(w, messResponse)
}

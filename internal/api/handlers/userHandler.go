package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/boginskiy/psychologistAI/internal/adapters/dto"
	"github.com/boginskiy/psychologistAI/internal/api"
	"github.com/boginskiy/psychologistAI/internal/api/handlers/mess"
	"github.com/boginskiy/psychologistAI/internal/api/response"
	"github.com/boginskiy/psychologistAI/internal/models"
	"github.com/boginskiy/psychologistAI/internal/service"
	"github.com/boginskiy/psychologistAI/internal/service/errs"
	"github.com/boginskiy/psychologistAI/pkg/request"
	"github.com/go-chi/chi"
)

type UserHandler struct {
	UserService service.UserService
	Sender      api.Sender
	basepath    string
}

func NewUserHandler(bpath string, userServ service.UserService, sender api.Sender) *UserHandler {
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

	token, err = h.UserService.Login(r.Context(), loginUser)

	// Errors / Убрать обработку ошибок в отдельную сущность
	if err != nil {
		switch {

		// Error Verification
		case errors.Is(err, errs.ErrVerification):
			messResponse.UpdateInfo(mess.MessNeedVerifyAccount, http.StatusForbidden)
		case errors.Is(err, errs.ErrAttemptsVerification):
			messResponse.UpdateInfo(mess.MessExceedingVerificationAttempts, http.StatusTooManyRequests)

		// Error

		default:
			messResponse.UpdateErr(err, http.StatusBadRequest)
		}
		h.Sender.SendResponse(w, messResponse)
		return
	}

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
	if err != nil {
		messResponse.UpdateErr(err, http.StatusBadRequest)
		h.Sender.SendResponse(w, messResponse)
		return
	}

	msg := "verification was successful"

	messResponse.UpdateInfo(msg, http.StatusOK)
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
		messResponse.UpdateErr(err, http.StatusBadRequest)
		h.Sender.SendResponse(w, messResponse)
		return
	}

	msg := fmt.Sprintf("go to '%s' and verify the account for %v minutes",
		userDomen.Email, models.TokenLifetime.Minutes())

	messResponse.UpdateInfo(msg, http.StatusOK)
	h.Sender.SendResponse(w, messResponse)
}

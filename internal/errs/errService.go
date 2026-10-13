package errs

import (
	"errors"
	"log"
	"net/http"

	"github.com/boginskiy/psychologistAI/internal/api"
	"github.com/boginskiy/psychologistAI/internal/api/response"
	"github.com/boginskiy/psychologistAI/internal/web/renders"
)

type ErrService struct {
	Responder api.Responder
}

// Доделываем responder, нужно чтобы вWEB было то что относится к ... а вAPIчто к ...

func NewErrService(responder api.Responder) *ErrService {
	return &ErrService{Responder: responder}
}

func (e *ErrService) CatchWebErr(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	// + logger
	log.Printf("web error from service: %v\n", err)

	switch {
	// Credentials
	case errors.Is(err, ErrInvalidCredentials):
		renders.RenderError(w, "401", http.StatusUnauthorized, nil)

	// Verification
	case errors.Is(err, ErrVerification):
		renders.RenderError(w, "403", http.StatusForbidden, nil)
	case errors.Is(err, ErrAttemptsVerification):
		renders.RenderError(w, "429", http.StatusTooManyRequests, nil)

	// Server
	case errors.Is(err, ErrServer):
		renders.RenderError(w, "500", http.StatusInternalServerError, nil)

	default:
		renders.RenderError(w, "400", http.StatusBadRequest, nil)
	}
}

func (e *ErrService) CatchApiErr(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	// + logger
	log.Printf("api error from service: %v\n", err)

	var body *response.InfoBody

	switch {
	// Credentials
	case errors.Is(err, ErrSession):
		body = response.NewInfoBodyWithErr(ErrAuth, http.StatusUnauthorized)
	case errors.Is(err, ErrUsingToken), errors.Is(err, ErrLegitimacyUser):
		body = response.NewInfoBodyWithErr(ErrAuth, http.StatusUnauthorized)
	case errors.Is(err, ErrTimeLiveSession), errors.Is(err, ErrCompareToken):
		body = response.NewInfoBodyWithErr(ErrAuth, http.StatusUnauthorized)

	// Server
	case errors.Is(err, ErrServer), errors.Is(err, ErrUpdateDBAfterRefresh):
		body = response.NewInfoBodyWithErr(ErrServer, http.StatusInternalServerError)

	default:
		body = response.NewInfoBodyWithErr(err, http.StatusBadRequest)
	}

	e.Responder.SendResponse(w, body)
}

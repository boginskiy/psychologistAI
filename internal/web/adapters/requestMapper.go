package adapters

import (
	"net/http"

	"github.com/boginskiy/psychologistAI/internal/adapters/dto"
	"github.com/boginskiy/psychologistAI/pkg/utils"
)

func ToCreateUserFromFormRequest(r *http.Request) (*dto.CreateUser, error) {
	if err := r.ParseForm(); err != nil {
		return nil, err
	}
	return &dto.CreateUser{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
		Name:     r.FormValue("name"),
		Phone:    r.FormValue("phone"),
		// Path for verification user
		VerificationLink: utils.Scheme(r) + "://" + r.Host + "/user/verification/",
	}, nil
}

func ToLoginUserFromFormRequest(r *http.Request) (*dto.LoginUser, error) {
	if err := r.ParseForm(); err != nil {
		return nil, err
	}
	return &dto.LoginUser{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}, nil
}

func ToMapRegistrTemplate(email string, liveTimeToken int) map[string]any {
	return map[string]any{
		"Email":         email,
		"LiveTimeToken": liveTimeToken,
	}
}

func ToMapStartTemplate(isUser bool) map[string]any {
	statusAuth := "LOG IN"
	if isUser {
		statusAuth = "LOG OUT"
	}

	return map[string]any{
		"StatusAuth": statusAuth,
	}
}

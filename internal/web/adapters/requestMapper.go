package adapters

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/boginskiy/psychologistAI/internal/adapters/dto"
	"github.com/boginskiy/psychologistAI/internal/api/request"
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

func ToInfoUserFromRequest(r *http.Request) (*dto.InfoUser, error) {
	infoUser, isUser := request.GetInfoUserFromContext(r.Context())

	if err := r.ParseForm(); err != nil {
		return nil, fmt.Errorf("error info user from request: %w", err)
	}

	msg := strings.TrimSpace(r.PostFormValue("message"))

	if !isUser {
		return &dto.InfoUser{Message: msg}, nil
	}

	infoUser.Message = msg
	return infoUser, nil
}

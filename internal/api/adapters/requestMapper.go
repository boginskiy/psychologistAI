package adapters

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/boginskiy/psychologistAI/cmd/config"
	"github.com/boginskiy/psychologistAI/internal/adapters/dto"
	"github.com/boginskiy/psychologistAI/internal/api"
	"github.com/boginskiy/psychologistAI/pkg/utils"
)

func NewTokenReq(requester api.Requester, request *http.Request, cookie *http.Cookie) *dto.TokenReq {
	os, browser, device := requester.TakeDeviceInfo(request)

	return &dto.TokenReq{
		Token:     cookie.Value,
		IP:        requester.TakeRealUserIP(request),
		UserAgent: requester.TakeUserAgent(request),
		OS:        os,
		Browser:   browser,
		Device:    device,
	}
}

func UpdateLoginUserFromRequest(loginUser *dto.LoginUser, requester api.Requester, request *http.Request) *dto.LoginUser {
	loginUser.OS, loginUser.Browser, loginUser.Device = requester.TakeDeviceInfo(request)
	loginUser.UserAgent = requester.TakeUserAgent(request)
	loginUser.IP = requester.TakeRealUserIP(request)
	return loginUser
}

func ToCreateUserFromRequest(r *http.Request) (*dto.CreateUser, error) {
	createUser := &dto.CreateUser{}
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read request body: %w", err)
	}

	if err := json.Unmarshal(body, &createUser); err != nil {
		return nil, fmt.Errorf("failed to deserialization request body: %w", err)
	}

	// Path for verification user
	createUser.VerificationLink = utils.Scheme(r) + "://" + r.Host + "/api/" + config.VersionAPI + "/user/verification/"

	return createUser, nil
}

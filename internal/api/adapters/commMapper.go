package adapters

import (
	"net/http"

	"github.com/boginskiy/psychologistAI/internal/adapters/dto"
	"github.com/boginskiy/psychologistAI/internal/api"
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

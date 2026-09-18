package api

import (
	"net/http"
)

type StatusGetter interface {
	GetStatus() int
}

type ErrorUpdater interface {
	ErrorUpdate(err error, status int)
}
type InfoUpdater interface {
	InfoUpdate(msg string, status int)
}

type ResponseBody interface {
	ErrorUpdater
	StatusGetter
	InfoUpdater
}

type ResponseSender interface {
	SetContentType(w http.ResponseWriter, contentType string)
	AddSetCookies(w http.ResponseWriter, cookies ...*http.Cookie)
	SendResponse(w http.ResponseWriter, body ResponseBody)
}

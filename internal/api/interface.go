package api

import "net/http"

// ResponseBody

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

// Responder

type Sender interface {
	SendResponse(w http.ResponseWriter, body ResponseBody)
}

type Setter interface {
	SetContentType(w http.ResponseWriter, contentType string)
}

type Adder interface {
	AddSetCookies(w http.ResponseWriter, cookies ...*http.Cookie)
}

type Responder interface {
	Setter
	Sender
	Adder
}

// Requester

type Reader interface {
	ReadAllRequestBody(r *http.Request, item any) (any, error)
}

type Taker interface {
	TakeRealUserIP(r *http.Request) string
	TakeUserAgent(r *http.Request) string
	TakeDeviceInfo(r *http.Request) (os, browser, device string)
}

type Requester interface {
	Reader
	Taker
}

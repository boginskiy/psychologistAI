package api

import (
	"net/http"
)

type ResponseWriter interface {
	GetStatus() int
	UpdateErr(err error, status int)
	UpdateInfo(info string, status int)
}

type Sender interface {
	SendResponse(w http.ResponseWriter, item ResponseWriter)
}

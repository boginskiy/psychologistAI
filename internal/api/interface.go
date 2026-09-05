package api

import (
	"net/http"
)

type Getter interface {
	GetStatus() int
}

type Preparator interface {
	PrepareOKResponse(status int)
	PrepareErrResponse(err error, status int)
}

type ResponseWriter interface {
	Getter
	Preparator
}

type Sender interface {
	SendResponse(w http.ResponseWriter, item ResponseWriter)
}

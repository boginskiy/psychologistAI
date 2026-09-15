package api

import (
	"net/http"

	models "github.com/boginskiy/psychologistAI/internal/models/users"
)

type ResponseWriter interface {
	GetStatus() int
}

type UpdaterResponse interface {
	ErrorUpdate(err error, status int)
	InfoUpdate(msg string, status int)
	AttrsUpdate(token *models.Token, status int)
}

type Sender interface {
	SendResponse(w http.ResponseWriter, item ResponseWriter)
}

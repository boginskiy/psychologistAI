package handlers

import "github.com/go-chi/chi"

type Registrar interface {
	Registration(r chi.Router)
}

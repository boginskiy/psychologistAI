package handlers

import "github.com/go-chi/chi"

type HomeHandler struct {
}

func (h *HomeHandler) Registration(r chi.Router) {
	r.Get("/", listUsers)   // GET /api/v1/user
	r.Post("/", createUser) // POST /api/v1/user
	r.Route("/{userID}", func(r chi.Router) {
		r.Get("/", getUser)    // GET /api/v1/user/{userID}
		r.Put("/", updateUser) // PUT /api/v1/user/{userID}

	})
}

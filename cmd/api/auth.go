package main

import (
	"net/http"

	"github.com/anvidev/project-time-tracker/internal/service/auth"
)

func (api *api) authRegister(w http.ResponseWriter, r *http.Request) {
	var body auth.RegisterRequest

	if err := api.readJSON(w, r, &body); err != nil {
		api.internalServerError(w, r, err)
		return
	}

	user, err := api.service.Auth.Register(r.Context(), body)
	if err != nil {
		switch err {
		case auth.ErrInvalidPassword:
			api.badRequestError(w, r, err)
		case auth.ErrDublicateEmail:
			api.conflictError(w, r, err)
		default:
			api.internalServerError(w, r, err)
		}
		return
	}

	response := map[string]any{
		"user": user,
	}

	if err := api.writeJSON(w, http.StatusCreated, response); err != nil {
		api.internalServerError(w, r, err)
		return
	}
}

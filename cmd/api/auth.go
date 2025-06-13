package main

import (
	"net/http"

	"github.com/anvidev/project-time-tracker/internal/store/users"
)

func (api *api) authRegister(w http.ResponseWriter, r *http.Request) {
	var body users.RegisterUserInput

	if err := api.readJSON(w, r, &body); err != nil {
		api.internalServerError(w, r, err)
		return
	}

	user, err := api.store.Users.Register(r.Context(), body)
	if err != nil {
		switch err {
		case users.ErrInvalidPassword:
			api.badRequestError(w, r, err)
		case users.ErrDublicateEmail:
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

func (api *api) authLogin(w http.ResponseWriter, r *http.Request) {
	var body users.LoginUserRequest

	if err := api.readJSON(w, r, &body); err != nil {
		api.internalServerError(w, r, err)
		return
	}

	ctx := r.Context()

	user, err := api.store.Users.GetByEmail(ctx, body.Email)
	if err != nil {
		switch err {
		case users.ErrUserNotFound:
			api.unauthorizedError(w, r, users.ErrInvalidCredentials)
		default:
			api.internalServerError(w, r, err)
		}
		return
	}

	if !user.IsActive {
		api.unauthorizedError(w, r, users.ErrUserNotActive)
		return
	}

	session, err := api.store.Sessions.Create(ctx, user.Id)
	if err != nil {
		api.internalServerError(w, r, err)
		return
	}

	response := map[string]any{
		"session": session,
	}

	if err := api.writeJSON(w, http.StatusCreated, response); err != nil {
		api.internalServerError(w, r, err)
		return
	}
}

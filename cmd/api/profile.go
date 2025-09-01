package main

import (
	"net/http"

	"github.com/anvidev/project-time-tracker/internal/store/users"
)

func (api *api) userProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId, _ := getUserId(ctx)

	user, err := api.store.Users.GetById(ctx, userId)
	if err != nil {
		switch err {
		case users.ErrUserNotFound:
			api.notFoundError(w, r, err)
		default:
			api.internalServerError(w, r, err)
		}
		return
	}

	response := map[string]any{
		"user": user,
	}

	if err := api.writeJSON(w, http.StatusOK, response); err != nil {
		api.internalServerError(w, r, err)
		return
	}
}

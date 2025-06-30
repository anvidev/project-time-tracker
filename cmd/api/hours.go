package main

import (
	"net/http"

	"github.com/anvidev/project-time-tracker/internal/store/hours"
)

func (api *api) hoursAll(w http.ResponseWriter, r *http.Request) {
	userId, _ := getUserId(r.Context())

	weekdays, err := api.store.Hours.AllWeekdays(r.Context(), userId)
	if err != nil {
		api.internalServerError(w, r, err)
		return
	}

	response := map[string]any{
		"hours": weekdays,
	}

	w.Header().Add("Cache-Control", "private, max-age=600")

	if err := api.writeJSON(w, http.StatusOK, response); err != nil {
		api.internalServerError(w, r, err)
		return
	}
}

func (api *api) update(w http.ResponseWriter, r *http.Request) {
	userId, _ := getUserId(r.Context())

	var body struct {
		Hours []hours.Weekday `json:"hours"`
	}
	if err := api.readJSON(w, r, &body); err != nil {
		api.badRequestError(w, r, err)
		return
	}

	err := api.store.Hours.UpdateWeekdays(r.Context(), userId, body.Hours)
	if err != nil {
		api.internalServerError(w, r, err)
		return
	}

	if err := api.writeJSON(w, http.StatusNoContent, nil); err != nil {
		api.internalServerError(w, r, err)
		return
	}
}

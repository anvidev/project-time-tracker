package main

import "net/http"

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

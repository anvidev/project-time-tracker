package main

import (
	"net/http"
	"time"

	"github.com/anvidev/project-time-tracker/internal/store/time_entries"
)

func (api *api) adminTimeEntries(w http.ResponseWriter, r *http.Request) {
	var filters time_entries.Filters

	if err := filters.Parse(r); err != nil {
		switch err {
		case
			time_entries.ErrInvalidCategoryId,
			time_entries.ErrInvalidUserId,
			time_entries.ErrInvalidFromDate,
			time_entries.ErrInvalidToDate,
			time_entries.ErrFromDateAfterToDate,
			time_entries.ErrToDateBeforeFromDate:
			api.badRequestError(w, r, err)
		default:
			api.internalServerError(w, r, err)
		}
		return
	}

	entries, err := api.store.TimeEntries.List(r.Context(), filters)
	if err != nil {
		api.internalServerError(w, r, err)
		return
	}

	var timeSpent time.Duration

	for _, entry := range entries {
		timeSpent += entry.Duration.Duration
	}

	response := map[string]any{
		"timeSpent": timeSpent.String(),
		"entries":   entries,
	}

	if err := api.writeJSON(w, http.StatusOK, response); err != nil {
		api.internalServerError(w, r, err)
		return
	}
}

func (api *api) adminUsers(w http.ResponseWriter, r *http.Request) {
	users, err := api.store.Users.List(r.Context())
	if err != nil {
		api.internalServerError(w, r, err)
		return
	}

	response := map[string]any{
		"user": users,
	}

	if err := api.writeJSON(w, http.StatusOK, response); err != nil {
		api.internalServerError(w, r, err)
		return
	}
}

func (api *api) adminCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := api.store.Categories.List(r.Context())
	if err != nil {
		api.internalServerError(w, r, err)
		return
	}

	response := map[string]any{
		"categories": categories,
	}

	if err := api.writeJSON(w, http.StatusOK, response); err != nil {
		api.internalServerError(w, r, err)
		return
	}
}

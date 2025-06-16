package main

import "net/http"

func (api *api) entriesCategories(w http.ResponseWriter, r *http.Request) {
	userId, _ := getUserId(r.Context())

	leafCategories, err := api.store.Categories.Leafs(r.Context(), userId)
	if err != nil {
		api.internalServerError(w, r, err)
		return
	}

	response := map[string]any{
		"categories": leafCategories,
	}

	w.Header().Add("Cache-Control", "private, max-age=600")

	if err := api.writeJSON(w, http.StatusOK, response); err != nil {
		api.internalServerError(w, r, err)
		return
	}
}

package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/anvidev/project-time-tracker/internal/store/categories"
	"github.com/anvidev/project-time-tracker/internal/store/time_entries"
)

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

func (api *api) entriesCategoriesTree(w http.ResponseWriter, r *http.Request) {
	userId, _ := getUserId(r.Context())

	categoriesTree, err := api.store.Categories.Tree(r.Context(), userId)
	if err != nil {
		api.internalServerError(w, r, err)
		return
	}

	response := map[string]any{
		"categories": categoriesTree,
	}

	w.Header().Add("Cache-Control", "private, max-age=600")

	if err := api.writeJSON(w, http.StatusOK, response); err != nil {
		api.internalServerError(w, r, err)
		return
	}
}

func (api *api) entriesRegisterTime(w http.ResponseWriter, r *http.Request) {
	userId, _ := getUserId(r.Context())

	var body time_entries.RegisterTimeEntryInput

	if err := api.readJSON(w, r, &body); err != nil {
		api.badRequestError(w, r, err)
		return
	}

	timeEntry, err := api.store.TimeEntries.Register(r.Context(), userId, body)
	if err != nil {
		api.internalServerError(w, r, err)
		return
	}

	response := map[string]any{
		"timeEntry": timeEntry,
	}

	if err := api.writeJSON(w, http.StatusCreated, response); err != nil {
		api.internalServerError(w, r, err)
		return
	}
}

func (api *api) entriesSummaryDay(w http.ResponseWriter, r *http.Request) {
	userId, _ := getUserId(r.Context())

	date := r.PathValue("date")
	_, err := time.Parse(time.DateOnly, date)
	if err != nil {
		api.badRequestError(w, r, err)
		return
	}

	summary, err := api.store.TimeEntries.SummaryDay(r.Context(), userId, date)
	if err != nil {
		api.internalServerError(w, r, err)
		return
	}

	response := map[string]any{
		"summary": summary,
	}

	if err := api.writeJSON(w, http.StatusOK, response); err != nil {
		api.internalServerError(w, r, err)
		return
	}
}

func (api *api) entriesSummaryMonth(w http.ResponseWriter, r *http.Request) {
	userId, _ := getUserId(r.Context())

	monthParts := strings.Split(r.PathValue("year-month"), "-")
	if len(monthParts) != 2 {
		api.badRequestError(w, r, fmt.Errorf("invalid month format"))
		return
	}

	year, err := strconv.Atoi(monthParts[0])
	if err != nil {
		api.badRequestError(w, r, err)
		return
	}

	month, err := strconv.ParseInt(monthParts[1], 10, 64)
	if err != nil {
		api.badRequestError(w, r, err)
		return
	}

	summary, err := api.store.TimeEntries.SummaryMonth(r.Context(), userId, time.Month(month), year)
	if err != nil {
		api.internalServerError(w, r, err)
		return
	}

	response := map[string]any{
		"summary": summary,
	}

	if err := api.writeJSON(w, http.StatusOK, response); err != nil {
		api.internalServerError(w, r, err)
		return
	}
}

func (api *api) entriesDelete(w http.ResponseWriter, r *http.Request) {
	userId, _ := getUserId(r.Context())

	entryId, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		api.badRequestError(w, r, err)
		return
	}

	if err := api.store.TimeEntries.Delete(r.Context(), entryId, userId); err != nil {
		api.internalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (api *api) entriesFollowCategory(w http.ResponseWriter, r *http.Request) {
	userId, _ := getUserId(r.Context())

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		api.badRequestError(w, r, err)
		return
	}

	if err := api.store.Categories.Follow(r.Context(), id, userId); err != nil {
		switch err {
		case categories.ErrAlreadyFollowed:
			api.conflictError(w, r, err)
		default:
			api.internalServerError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (api *api) entriesUnfollowCategory(w http.ResponseWriter, r *http.Request) {
	userId, _ := getUserId(r.Context())

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		api.badRequestError(w, r, err)
		return
	}

	if err := api.store.Categories.Unfollow(r.Context(), id, userId); err != nil {
		switch err {
		case categories.ErrNotFollowingCategory:
			api.notFoundError(w, r, err)
		default:
			api.internalServerError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

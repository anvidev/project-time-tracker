package main

import (
	"errors"
	"io"
	"net/http"

	"github.com/anvidev/project-time-tracker/internal/store/live_timers"
)

func (api *api) timerError(w http.ResponseWriter, r *http.Request, err error) {
	switch err {
	case live_timers.ErrTimerNotFound:
		api.notFoundError(w, r, err)
	case live_timers.ErrTimerExists,
		live_timers.ErrTimerNotRunning,
		live_timers.ErrTimerNotPaused:
		api.conflictError(w, r, err)
	case live_timers.ErrNothingToSave, live_timers.ErrInvalidDate:
		api.badRequestError(w, r, err)
	default:
		api.internalServerError(w, r, err)
	}
}

func (api *api) timerGet(w http.ResponseWriter, r *http.Request) {
	userId, _ := getUserId(r.Context())

	timer, err := api.store.LiveTimers.Get(r.Context(), userId)
	if err != nil && err != live_timers.ErrTimerNotFound {
		api.internalServerError(w, r, err)
		return
	}

	// no active timer is a valid state, so clients get timer: null instead of a 404
	response := map[string]any{
		"timer": timer,
	}

	if err := api.writeJSON(w, http.StatusOK, response); err != nil {
		api.internalServerError(w, r, err)
		return
	}
}

func (api *api) timerStart(w http.ResponseWriter, r *http.Request) {
	userId, _ := getUserId(r.Context())

	var body live_timers.StartTimerInput

	if err := api.readJSON(w, r, &body); err != nil {
		api.badRequestError(w, r, err)
		return
	}

	timer, err := api.store.LiveTimers.Start(r.Context(), userId, body)
	if err != nil {
		api.timerError(w, r, err)
		return
	}

	if err := api.writeJSON(w, http.StatusCreated, map[string]any{"timer": timer}); err != nil {
		api.internalServerError(w, r, err)
		return
	}
}

func (api *api) timerPause(w http.ResponseWriter, r *http.Request) {
	userId, _ := getUserId(r.Context())

	timer, err := api.store.LiveTimers.Pause(r.Context(), userId)
	if err != nil {
		api.timerError(w, r, err)
		return
	}

	if err := api.writeJSON(w, http.StatusOK, map[string]any{"timer": timer}); err != nil {
		api.internalServerError(w, r, err)
		return
	}
}

func (api *api) timerResume(w http.ResponseWriter, r *http.Request) {
	userId, _ := getUserId(r.Context())

	timer, err := api.store.LiveTimers.Resume(r.Context(), userId)
	if err != nil {
		api.timerError(w, r, err)
		return
	}

	if err := api.writeJSON(w, http.StatusOK, map[string]any{"timer": timer}); err != nil {
		api.internalServerError(w, r, err)
		return
	}
}

func (api *api) timerUpdate(w http.ResponseWriter, r *http.Request) {
	userId, _ := getUserId(r.Context())

	var body live_timers.UpdateTimerInput

	if err := api.readJSON(w, r, &body); err != nil {
		api.badRequestError(w, r, err)
		return
	}

	timer, err := api.store.LiveTimers.Update(r.Context(), userId, body)
	if err != nil {
		api.timerError(w, r, err)
		return
	}

	if err := api.writeJSON(w, http.StatusOK, map[string]any{"timer": timer}); err != nil {
		api.internalServerError(w, r, err)
		return
	}
}

func (api *api) timerSave(w http.ResponseWriter, r *http.Request) {
	userId, _ := getUserId(r.Context())

	var body live_timers.SaveTimerInput

	// the body is optional, so an empty request saves the tracked time as is
	if err := api.readJSON(w, r, &body); err != nil && !errors.Is(err, io.EOF) {
		api.badRequestError(w, r, err)
		return
	}

	entry, err := api.store.LiveTimers.Save(r.Context(), userId, body)
	if err != nil {
		api.timerError(w, r, err)
		return
	}

	if err := api.writeJSON(w, http.StatusCreated, map[string]any{"timeEntry": entry}); err != nil {
		api.internalServerError(w, r, err)
		return
	}
}

func (api *api) timerDiscard(w http.ResponseWriter, r *http.Request) {
	userId, _ := getUserId(r.Context())

	if err := api.store.LiveTimers.Discard(r.Context(), userId); err != nil {
		api.timerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

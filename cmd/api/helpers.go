package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

const (
	maxBodySize int64 = 1_048_576 // 1mb
)

var (
	validate *validator.Validate

	ErrorCodeInternal        string = "INTERNAL_SERVER_ERROR"
	ErrorCodeBadRequest             = "BAD_REQUEST"
	ErrorCodeNotFound               = "NOT_FOUND"
	ErrorCodeConflict               = "CONFLICT"
	ErrorCodeUnauthorized           = "UNAUTHORIZED"
	ErrorCodeTooManyRequests        = "TOO_MANY_REQUESTS"
)

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}

type errorEnvelope struct {
	Error string `json:"error"`
	Code  string `json:"code"`
}

func newErrorEnvelope(errMsg, code string) errorEnvelope {
	return errorEnvelope{Error: errMsg, Code: code}
}

// json helpers
func (api *api) readJSON(w http.ResponseWriter, r *http.Request, v any) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(v); err != nil {
		return fmt.Errorf("invalid json: %w", err)
	}
	return validate.Struct(v)
}

func (api *api) writeJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

// error response helpers
func (api *api) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	api.logger.Error("internal server error",
		"method", r.Method,
		"path", r.URL.Path,
		"error", err.Error())

	api.writeJSON(w, http.StatusInternalServerError, newErrorEnvelope("something went wrong", ErrorCodeInternal))
}

func (api *api) badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	api.logger.Info("bad request error",
		"method", r.Method,
		"path", r.URL.Path,
		"error", err.Error())

	api.writeJSON(w, http.StatusBadRequest, newErrorEnvelope(err.Error(), ErrorCodeBadRequest))
}

func (api *api) notFoundError(w http.ResponseWriter, r *http.Request, err error) {
	api.logger.Warn("not found error",
		"method", r.Method,
		"path", r.URL.Path,
		"error", err.Error())

	api.writeJSON(w, http.StatusNotFound, newErrorEnvelope(err.Error(), ErrorCodeNotFound))
}

func (api *api) conflictError(w http.ResponseWriter, r *http.Request, err error) {
	api.logger.Error("conflict error",
		"method", r.Method,
		"path", r.URL.Path,
		"error", err.Error())

	api.writeJSON(w, http.StatusConflict, newErrorEnvelope(err.Error(), ErrorCodeConflict))
}

func (api *api) unauthorizedError(w http.ResponseWriter, r *http.Request, err error) {
	api.logger.Warn("unauthorized error",
		"method", r.Method,
		"path", r.URL.Path,
		"error", err.Error())

	api.writeJSON(w, http.StatusUnauthorized, newErrorEnvelope(err.Error(), ErrorCodeUnauthorized))
}

func (api *api) tooManyRequestsError(w http.ResponseWriter, r *http.Request, err error) {
	api.logger.Warn("too many requests error",
		"method", r.Method,
		"path", r.URL.Path,
		"error", err.Error())

	api.writeJSON(w, http.StatusTooManyRequests, newErrorEnvelope(err.Error(), ErrorCodeTooManyRequests))
}

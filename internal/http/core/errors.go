package core

import (
	"encoding/json"
	"errors"
	"net/http"
)

var (
	ErrMethodNotAllowed = errors.New("method not allowed")
	ErrNotFound         = errors.New("resource not found")
	ErrInvalidInput     = errors.New("invalid input")
	ErrConflict         = errors.New("resource conflict")
	ErrInternal         = errors.New("internal server error")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrForbidden        = errors.New("forbidden")
	ErrBadRequest       = errors.New("bad request")
)

type ApiError struct {
	Code    int
	Message string
	Err     error
}

func (e *ApiError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}

	return e.Message
}

func (e ApiError) WriteHTTP(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.Code)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    e.Code,
		"message": e.Message,
	})
}

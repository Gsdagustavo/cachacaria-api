package core

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

var (
	// HTTP
	ErrMethodNotAllowed = errors.New("method not allowed")
	ErrNotFound         = errors.New("resource not found")
	ErrInvalidInput     = errors.New("invalid input")
	ErrConflict         = errors.New("resource conflict")
	ErrInternal         = errors.New("internal server error")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrForbidden        = errors.New("forbidden")
	ErrBadRequest       = errors.New("bad request")

	// Product
	ErrInvalidProductName  = errors.New("invalid product name")
	ErrInvalidProductPrice = errors.New("invalid product price")
	ErrInvalidProductStock = errors.New("invalid product stock")
	ErrNoProductPhoto      = errors.New("no product photo")

	// User
	ErrInvalidEmail       = errors.New("invalid email")
	ErrInvalidPassword    = errors.New("invalid password")
	ErrInvalidPhoneNumber = errors.New("invalid phone number")
	ErrInvalidToken       = errors.New("invalid token")
	ErrTokenExpired       = errors.New("token expired")
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

func (e *ApiError) WithError(location string) *ApiError {
	if e.Err != nil {
		log.Printf("[ERROR] on %s: %s | cause: %v", location, e.Message, e.Err)
	} else {
		log.Printf("[ERROR] on %s - %s", location, e.Message)
	}

	return e
}

func (e *ApiError) WriteHTTP(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.Code)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    e.Code,
		"message": e.Message,
	})
}

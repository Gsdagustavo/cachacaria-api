package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var (
	// HTTP
	ErrMethodNotAllowed = errors.New("method not allowed")
	ErrNotFound         = errors.New("resource not found")
	ErrInvalidInput     = errors.New("invalid input")
	ErrConflict         = errors.New("resource conflict222")
	ErrInternal         = errors.New("Erro interno no servidor")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrForbidden        = errors.New("forbidden")
	ErrBadRequest       = errors.New("bad request")

	// Product
	ErrInvalidProductName  = errors.New("invalid product name")
	ErrInvalidProductPrice = errors.New("invalid product price")
	ErrInvalidProductStock = errors.New("invalid product stock")
	ErrNoProductPhoto      = errors.New("no product photo")

	// User
	ErrTokenGenerationError = errors.New("token generation error")
	ErrUserAlreadyExists    = errors.New("user already exists")
	ErrUserNotFound         = errors.New("user not found")
	ErrInvalidEmail         = errors.New("invalid email")
	ErrInvalidPassword      = errors.New("invalid password")
	ErrInvalidPhoneNumber   = errors.New("invalid phone number")
	ErrInvalidToken         = errors.New("invalid token")
	ErrTokenExpired         = errors.New("token expired")
)

type ServerError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Cause   string `json:"cause,omitempty"`
	Err     error  `json:"-"`
}

func (e *ServerError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[HTTP %d] %s | cause: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[HTTP %d] %s", e.Code, e.Message)
}

func (e *ServerError) Unwrap() error {
	return e.Err
}

func (e *ServerError) WithError(context string) *ServerError {
	if e == nil {
		return nil
	}
	e.Cause = context
	return e
}

func (e *ServerError) WriteHTTP(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.Code)

	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    e.Code,
		"message": e.Message,
	}); err != nil {
		fmt.Printf("failed to encode error response: %v\n", err)
	}
}

// === Constructors (HTTP) ===
func BadRequest(msg string, err error) *ServerError {
	if msg == "" {
		msg = ErrBadRequest.Error()
	}
	return &ServerError{Code: http.StatusBadRequest, Message: msg, Err: err}
}

func Unauthorized(msg string, err error) *ServerError {
	if msg == "" {
		msg = ErrUnauthorized.Error()
	}
	return &ServerError{Code: http.StatusUnauthorized, Message: msg, Err: err}
}

func Forbidden(msg string, err error) *ServerError {
	if msg == "" {
		msg = ErrForbidden.Error()
	}
	return &ServerError{Code: http.StatusForbidden, Message: msg, Err: err}
}

func NotFound(msg string, err error) *ServerError {
	if msg == "" {
		msg = ErrNotFound.Error()
	}
	return &ServerError{Code: http.StatusNotFound, Message: msg, Err: err}
}

func MethodNotAllowed(msg string, err error) *ServerError {
	if msg == "" {
		msg = ErrMethodNotAllowed.Error()
	}
	return &ServerError{Code: http.StatusMethodNotAllowed, Message: msg, Err: err}
}

func Conflict(msg string, err error) *ServerError {
	if msg == "" {
		msg = ErrConflict.Error()
	}
	return &ServerError{Code: http.StatusConflict, Message: msg, Err: err}
}

func UnprocessableEntity(msg string, err error) *ServerError {
	if msg == "" {
		msg = ErrInvalidInput.Error()
	}
	return &ServerError{Code: http.StatusUnprocessableEntity, Message: msg, Err: err}
}

func Internal(msg string, err error) *ServerError {
	if msg == "" {
		msg = ErrInternal.Error()
	}
	return &ServerError{Code: http.StatusInternalServerError, Message: msg, Err: err}
}

// === Constructors (Domain-specific helpers) ===
// Product
func InvalidProductName(err error) *ServerError {
	return BadRequest(ErrInvalidProductName.Error(), err)
}
func InvalidProductPrice(err error) *ServerError {
	return BadRequest(ErrInvalidProductPrice.Error(), err)
}
func InvalidProductStock(err error) *ServerError {
	return BadRequest(ErrInvalidProductStock.Error(), err)
}
func NoProductPhoto(err error) *ServerError {
	return BadRequest(ErrNoProductPhoto.Error(), err)
}

// User/Auth
func UserAlreadyExists(err error) *ServerError {
	return Conflict(ErrUserAlreadyExists.Error(), err)
}
func UserNotFound(err error) *ServerError {
	return NotFound(ErrUserNotFound.Error(), err)
}
func InvalidEmail(err error) *ServerError {
	return BadRequest(ErrInvalidEmail.Error(), err)
}
func InvalidPassword(err error) *ServerError {
	return BadRequest(ErrInvalidPassword.Error(), err)
}
func InvalidPhoneNumber(err error) *ServerError {
	return BadRequest(ErrInvalidPhoneNumber.Error(), err)
}
func TokenGenerationError(err error) *ServerError {
	return Internal(ErrTokenGenerationError.Error(), err)
}
func InvalidToken(err error) *ServerError {
	return Unauthorized(ErrInvalidToken.Error(), err)
}
func TokenExpired(err error) *ServerError {
	return Unauthorized(ErrTokenExpired.Error(), err)
}

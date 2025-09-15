package authhandler

import (
	"cachacariaapi/internal/domain/entities"
	"cachacariaapi/internal/domain/usecases/user"
	"cachacariaapi/internal/interfaces/http/core"

	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	UserUseCases *userusecases.UserUseCases
	jtwSecret    []byte
}

func NewAuthHandler(userUseCases *userusecases.UserUseCases) *AuthHandler {
	secret := os.Getenv("JWT_SECRET")
	return &AuthHandler{UserUseCases: userUseCases, jtwSecret: []byte(secret)}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) *core.ApiError {
	if apiErr := core.ValidateRequestMethod(r, http.MethodPost); apiErr != nil {
		return apiErr.WithError("auth handler / register")
	}

	var request entities.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return (&core.ApiError{
			Code:    http.StatusBadRequest,
			Message: core.ErrBadRequest.Error(),
			Err:     err,
		}).WithError("auth handler / register")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return (&core.ApiError{
			Code:    http.StatusInternalServerError,
			Message: core.ErrInternal.Error(),
			Err:     err,
		}).WithError("auth handler / register")
	}
	request.Password = string(hashedPassword)

	user, err := h.UserUseCases.Add(request)
	if err != nil {
		return (MapRegisterError(err)).WithError("auth handler / register")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(h.jtwSecret)
	if err != nil {
		return (&core.ApiError{
			Code:    http.StatusInternalServerError,
			Message: core.ErrTokenGenerationError.Error(),
			Err:     err,
		}).WithError("auth handler / register")
	}

	setToken(w, tokenString)
	return nil
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) *core.ApiError {
	if apiErr := core.ValidateRequestMethod(r, http.MethodPost); apiErr != nil {
		return apiErr.WithError("auth handler / login")
	}

	var req entities.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return (&core.ApiError{
			Code:    http.StatusBadRequest,
			Message: core.ErrBadRequest.Error(),
			Err:     err,
		}).WithError("auth handler / login")
	}

	user, err := h.UserUseCases.FindByEmail(req.Email)
	if err != nil {
		return (MapLoginError(err)).WithError("auth handler / register")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return (&core.ApiError{
			Code:    http.StatusUnauthorized,
			Message: core.ErrInvalidPassword.Error(),
			Err:     core.ErrUnauthorized,
		}).WithError("auth handler / login")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(h.jtwSecret)
	if err != nil {
		return (&core.ApiError{
			Code:    http.StatusInternalServerError,
			Message: "token generation error",
			Err:     err,
		}).WithError("auth handler / login")
	}

	setToken(w, tokenString)
	return nil
}

func setToken(w http.ResponseWriter, token string) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Authorization", "Bearer "+token)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func MapRegisterError(err error) *core.ApiError {
	switch {
	case errors.Is(err, core.ErrConflict):
		return &core.ApiError{Code: http.StatusConflict, Message: core.ErrUserAlreadyExists.Error(), Err: err}
	case errors.Is(err, core.ErrInvalidEmail):
		return &core.ApiError{Code: http.StatusBadRequest, Message: core.ErrInvalidEmail.Error(), Err: err}
	case errors.Is(err, core.ErrInvalidPassword):
		return &core.ApiError{Code: http.StatusBadRequest, Message: core.ErrInvalidPassword.Error(), Err: err}
	case errors.Is(err, core.ErrInvalidPhoneNumber):
		return &core.ApiError{Code: http.StatusBadRequest, Message: core.ErrInvalidPhoneNumber.Error(), Err: err}
	default:
		return &core.ApiError{Code: http.StatusInternalServerError, Message: core.ErrInternal.Error(), Err: err}
	}
}

func MapLoginError(err error) *core.ApiError {
	switch {
	case errors.Is(err, core.ErrNotFound):
		return &core.ApiError{Code: http.StatusNotFound, Message: core.ErrUserNotfound.Error(), Err: err}
	case errors.Is(err, core.ErrInvalidPassword):
		return &core.ApiError{Code: http.StatusBadRequest, Message: core.ErrInvalidPassword.Error(), Err: err}
	default:
		return &core.ApiError{Code: http.StatusInternalServerError, Message: core.ErrInternal.Error(), Err: err}
	}
}

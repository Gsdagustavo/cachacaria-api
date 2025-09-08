package authhandler

import (
	"cachacariaapi/internal/domain/entities"
	"cachacariaapi/internal/domain/usecases"
	core2 "cachacariaapi/internal/interfaces/http/core"

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

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) *core2.ApiError {
	if apiErr := core2.ValidateRequestMethod(r, http.MethodPost); apiErr != nil {
		return apiErr.WithError("auth handler / register")
	}

	var request entities.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return (&core2.ApiError{
			Code:    http.StatusBadRequest,
			Message: "bad request",
			Err:     err,
		}).WithError("auth handler / register")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return (&core2.ApiError{
			Code:    http.StatusInternalServerError,
			Message: "password hashing error",
			Err:     err,
		}).WithError("auth handler / register")
	}
	request.Password = string(hashedPassword)

	user, err := h.UserUseCases.Add(request)
	if err != nil {
		if errors.Is(err, core2.ErrConflict) {
			return (&core2.ApiError{
				Code:    http.StatusConflict,
				Message: "user already exists",
				Err:     err,
			}).WithError("auth handler / register")
		}

		if errors.Is(err, core2.ErrBadRequest) {
			return (&core2.ApiError{
				Code:    http.StatusBadRequest,
				Message: "bad request",
				Err:     err,
			}).WithError("auth handler / register")
		}

		return (&core2.ApiError{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
			Err:     err,
		}).WithError("auth handler / register")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(h.jtwSecret)
	if err != nil {
		return (&core2.ApiError{
			Code:    http.StatusInternalServerError,
			Message: "token generation error",
			Err:     err,
		}).WithError("auth handler / register")
	}

	setToken(w, tokenString)
	return nil
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) *core2.ApiError {
	if apiErr := core2.ValidateRequestMethod(r, http.MethodPost); apiErr != nil {
		return apiErr.WithError("auth handler / login")
	}

	var req entities.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return (&core2.ApiError{
			Code:    http.StatusBadRequest,
			Message: "bad request",
			Err:     err,
		}).WithError("auth handler / login")
	}

	user, err := h.UserUseCases.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, core2.ErrNotFound) {
			return (&core2.ApiError{
				Code:    http.StatusNotFound,
				Message: "user not found",
				Err:     err,
			}).WithError("auth handler / login")
		}

		return (&core2.ApiError{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
			Err:     err,
		}).WithError("auth handler / login")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return (&core2.ApiError{
			Code:    http.StatusUnauthorized,
			Message: "invalid password",
			Err:     core2.ErrUnauthorized,
		}).WithError("auth handler / login")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(h.jtwSecret)
	if err != nil {
		return (&core2.ApiError{
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

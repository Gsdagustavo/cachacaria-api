package authhandler

import (
	"cachacariaapi/internal/http/core"
	"cachacariaapi/internal/models"
	"cachacariaapi/internal/usecases"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	UserUseCases usecases.UserUseCases
	jtwSecret    []byte
}

func NewAuthHandler(userUseCases usecases.UserUseCases) *AuthHandler {
	secret := os.Getenv("JWT_SECRET")
	return &AuthHandler{UserUseCases: userUseCases, jtwSecret: []byte(secret)}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) *core.ApiError {
	if apiErr := core.ValidateRequestMethod(r, http.MethodPost); apiErr != nil {
		return apiErr
	}

	var request models.RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Printf("err: %v", err)

		return &core.ApiError{
			Code:    http.StatusBadRequest,
			Message: "bad request",
			Err:     err,
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	request.Password = string(hashedPassword)
	user, err := h.UserUseCases.Add(request)
	if err != nil {
		if errors.Is(err, core.ErrConflict) {
			return &core.ApiError{
				Code:    http.StatusConflict,
				Message: "user already exists",
				Err:     err,
			}
		}

		return &core.ApiError{
			Code:    http.StatusBadRequest,
			Message: "bad request",
			Err:     err,
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(h.jtwSecret)
	if err != nil {
		return &core.ApiError{
			Code:    http.StatusInternalServerError,
			Message: "token generation error",
			Err:     err,
		}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    tokenString,
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: true,
		//Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
	})

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Authorization", "Bearer "+tokenString)
	json.NewEncoder(w).Encode(tokenString)
	return nil
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) *core.ApiError {
	if apiErr := core.ValidateRequestMethod(r, http.MethodPost); apiErr != nil {
		return apiErr
	}

	var req models.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("err: %v", err)

		return &core.ApiError{
			Code:    http.StatusBadRequest,
			Message: "bad request",
			Err:     err,
		}
	}

	user, err := h.UserUseCases.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, core.ErrNotFound) {
			return &core.ApiError{
				Code:    http.StatusNotFound,
				Message: "user not found",
				Err:     err,
			}
		}

		return &core.ApiError{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
			Err:     err,
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		log.Printf("err: %v", err)

		return &core.ApiError{
			Code:    http.StatusUnauthorized,
			Message: "invalid password",
			Err:     core.ErrUnauthorized,
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(h.jtwSecret)
	if err != nil {
		return &core.ApiError{
			Code:    http.StatusInternalServerError,
			Message: "token generation error",
			Err:     err,
		}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    tokenString,
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: true,
		//Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
	})

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Authorization", "Bearer "+tokenString)
	json.NewEncoder(w).Encode(tokenString)
	return nil
}

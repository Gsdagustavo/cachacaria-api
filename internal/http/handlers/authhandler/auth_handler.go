package authhandler

import (
	"cachacariaapi/internal/http/core"
	"cachacariaapi/internal/models"
	"cachacariaapi/internal/usecases"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

type AuthHandler struct {
	UserUseCases usecases.UserUseCases
}

func NewAuthHandler(userUseCases usecases.UserUseCases) *AuthHandler {
	return &AuthHandler{UserUseCases: userUseCases}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) *core.ApiError {
	if apiErr := core.ValidateRequestMethod(r, http.MethodPost); apiErr != nil {
		return apiErr
	}

	log.Println("Registering user")

	var user models.UserRequest
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		return &core.ApiError{
			Err: err,
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return &core.ApiError{
			Err: err,
		}
	}

	user.Password = string(hashedPassword)

	response, err := h.UserUseCases.Add(user)
	if err != nil {
		return &core.ApiError{
			Err: err,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	return nil
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) *core.ApiError {
	if apiErr := core.ValidateRequestMethod(r, http.MethodPost); apiErr != nil {
		return apiErr
	}

	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		return &core.ApiError{
			Err: err,
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		return &core.ApiError{
			Err: err,
		}
	}

	creds.Password = string(hashedPassword)

	user, err := h.UserUseCases.FindByEmail(creds.Email)
	if err != nil {
		return &core.ApiError{
			Err: err,
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(token)
	return nil
}

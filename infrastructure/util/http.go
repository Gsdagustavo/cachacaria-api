package util

import (
	"cachacariaapi/domain/entities"
	"context"
	"encoding/json"
	"log"
	"net/http"
)

// ContextKey is a private type used for context keys.
type ContextKey string

// UserIDContextKey is the key for the user ID in the request context.
const UserIDContextKey ContextKey = "userID"

func NewContextWithUserID(ctx context.Context, userID int) context.Context {
	return context.WithValue(ctx, UserIDContextKey, userID)
}

func GetUserIDFromContext(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value(UserIDContextKey).(int)
	return userID, ok
}

func ValidateRequestMethod(r *http.Request, allowedMethod string) *entities.ServerError {
	if r.Method != allowedMethod {
		return &entities.ServerError{
			Code: http.StatusMethodNotAllowed,
			Err:  nil,
		}
	}
	return nil
}

func WriteGenericResponse(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func WriteServerResponse(w http.ResponseWriter, response *entities.ServerResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.Code)

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func WriteServerError(w http.ResponseWriter, error *entities.ServerError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(error.Code)

	err := json.NewEncoder(w).Encode(error)
	if err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

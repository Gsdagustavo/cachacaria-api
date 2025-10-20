package util

import (
	"cachacariaapi/domain/entities"
	"context"
	"encoding/json"
	"io/ioutil"
	"log/slog"
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

func WriteResponse(w http.ResponseWriter, statusCode int, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	bytes, err := json.Marshal(response)
	if err != nil {
		slog.Error("error writing response", "cause", err)
		return
	}

	w.Write(bytes)
}

type ServerResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func Write(w http.ResponseWriter, v any) {
	bytes, err := json.Marshal(v)
	if err != nil {
		slog.Error("failed to marshal response", "cause", err)
		WriteInternalError(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(bytes)
	if err != nil {
		slog.Error("failed to write response", "cause", err)
		WriteInternalError(w)
	}
}

func WriteInternalError(w http.ResponseWriter) {
	http.Error(w, "Internal server error", http.StatusInternalServerError)
}

func WriteBadRequest(w http.ResponseWriter) {
	http.Error(w, "Bad request", http.StatusBadRequest)
}

func Read(r *http.Request, v any) error {
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, v)
	return nil
}

package util

import (
	"cachacariaapi/domain/entities"
	"context"
	"encoding/json"
	"io/ioutil"
	"log/slog"
	"net/http"
)

type ContextKey string

const userContextKey ContextKey = "userContext"

type UserContext struct {
	UserID  int
	IsAdmin bool
}

func NewContextWithUser(ctx context.Context, userID int, isAdmin bool) context.Context {
	return context.WithValue(ctx, userContextKey, &UserContext{
		UserID:  userID,
		IsAdmin: isAdmin,
	})
}

func GetUserFromContext(ctx context.Context) (*UserContext, bool) {
	user, ok := ctx.Value(userContextKey).(*UserContext)
	return user, ok
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

func WriteResponse(w http.ResponseWriter, response ServerResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.Status)
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

package userhandler

import (
	"cachacariaapi/internal/http/core"
	"cachacariaapi/internal/usecases"
	"encoding/json"
	"errors"
	"net/http"
)

type UserHandler struct {
	UserUseCases usecases.UserUseCases
}

func NewUserHandler(u usecases.UserUseCases) *UserHandler {
	return &UserHandler{u}
}

func ValidateRequestMethod(r *http.Request, allowedMethod string) *core.ApiError {
	if r.Method != allowedMethod {
		return &core.ApiError{
			Code:    http.StatusMethodNotAllowed,
			Message: core.ErrMethodNotAllowed.Error(),
			Err:     nil,
		}
	}
	return nil
}

func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	if apiErr := ValidateRequestMethod(r, http.MethodGet); apiErr != nil {
		apiErr.WriteHTTP(w)
		return
	}

	users, err := h.UserUseCases.GetAll()

	if err != nil {
		if errors.Is(err, core.ErrNotFound) {
			core.ApiError{
				Code:    http.StatusNotFound,
				Message: err.Error(),
			}.WriteHTTP(w)
			return
		}

		core.ApiError{
			Code:    http.StatusInternalServerError,
			Message: "could not get users",
		}.WriteHTTP(w)
	}

	json.NewEncoder(w).Encode(users)
}

//func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
//
//	if !ValidateRequestMethod(r, http.MethodGet, w) {
//		return
//	}
//
//	query := r.URL.Query()
//	idStr := query.Get("id")
//	if idStr == "" {
//
//		return
//	}
//
//	id, err := strconv.ParseInt(idStr, 10, 64)
//	if err != nil {
//		WriteJSONError(w, &core.ApiError{Code: 400, Message: "invalid id"})
//		return
//	}
//
//	user, err := h.UserUseCases.FindById(id)
//	if err != nil {
//		WriteJSONError(w, err)
//		return
//	}
//
//	w.Header().Set("Content-Type", "application/json")
//	json.NewEncoder(w).Encode(user)
//}
//
//func (h *UserHandler) AddUser(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("Content-Type", "application/json")
//
//	if !ValidateRequestMethod(r, http.MethodPost, w) {
//		return
//	}
//
//	var req models.UserRequest
//
//	err := json.NewDecoder(r.Body).Decode(&req)
//
//	if err != nil {
//		w.WriteHeader(http.StatusBadRequest)
//		log.Printf("Error adding user: %v", err)
//		return
//	}
//
//	res, err := h.UserUseCases.Add(req)
//
//	log.Printf("Added user: %v", req)
//
//	if err != nil {
//		w.WriteHeader(http.StatusBadRequest)
//		json.NewEncoder(w).Encode(core.ErrorMessage{Message: err.Error()})
//		return
//	}
//
//	prettyJSON, err := json.MarshalIndent(res, "", "  ")
//	if err != nil {
//		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
//		return
//	}
//
//	w.Write(prettyJSON)
//
//	log.Printf("Response: %v", res)
//}
//
//func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("Content-Type", "application/json")
//
//	if !ValidateRequestMethod(r, http.MethodDelete, w) {
//		return
//	}
//
//	query := r.URL.Query()
//
//	if !query.Has("id") {
//		w.WriteHeader(http.StatusBadRequest)
//		w.Write([]byte(`{"message":"id is required"}`))
//		return
//	}
//
//	id, err := strconv.ParseInt(query.Get("id"), 10, 64)
//
//	if err != nil {
//		w.WriteHeader(http.StatusBadRequest)
//		w.Write([]byte(`{"message":"invalid id"}`))
//		return
//	}
//
//	err = h.UserUseCases.Delete(id)
//	if err != nil {
//		w.WriteHeader(http.StatusNotFound)
//		log.Printf("User not found for the ID %v", id)
//		return
//	}
//
//	w.WriteHeader(http.StatusOK)
//	json.NewEncoder(w).Encode("User deleted")
//}

//func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("Content-Type", "application/json")
//
//	if r.Method != http.MethodPut {
//		w.WriteHeader(http.StatusMethodNotAllowed)
//		return
//	}
//
//	h.UserUseCases.Update()
//}

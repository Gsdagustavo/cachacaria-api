package userhandler

import (
	"cachacariaapi/internal/http/core"
	"cachacariaapi/internal/models"
	"cachacariaapi/internal/usecases"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type UserHandler struct {
	UserUseCases usecases.UserUseCases
}

func NewUserHandler(u usecases.UserUseCases) *UserHandler {
	return &UserHandler{u}
}

func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	users, err := h.UserUseCases.GetAll()

	log.Printf("Retrieved %d users", len(users))

	if err != nil {
		core.ErrorMessage{Code: http.StatusInternalServerError, Message: err.Error()}.ShowErrorMessage(w)
	}

	if len(users) == 0 {
		core.ErrorMessage{Code: http.StatusNotFound, Message: "No users found"}.ShowErrorMessage(w)
		return
	}

	json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query()

	if !query.Has("id") {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(core.ErrorMessage{Message: "id is required"})
		return
	}

	id, err := strconv.ParseInt(query.Get("id"), 10, 64)

	log.Printf("Query ID: %v", id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("User not found for the ID %v", id)
		return
	}

	user, err := h.UserUseCases.FindById(id)

	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		log.Printf("User not found for the ID %v", id)
		json.NewEncoder(w).Encode(core.ErrorMessage{Message: ""})
		return
	}

	prettyJSON, err := json.MarshalIndent(user, "", "  ")
	if err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(prettyJSON)

	fmt.Printf("User: %v", user)
}

func (h *UserHandler) AddUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req models.UserRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error adding user: %v", err)
		return
	}

	res, err := h.UserUseCases.Add(req)

	log.Printf("Added user: %v", req)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(core.ErrorMessage{Message: err.Error()})
		return
	}

	prettyJSON, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}

	w.Write(prettyJSON)

	log.Printf("Response: %v", res)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query()

	if !query.Has("id") {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message":"id is required"}`))
		return
	}

	id, err := strconv.ParseInt(query.Get("id"), 10, 64)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message":"invalid id"}`))
		return
	}

	err = h.UserUseCases.Delete(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Printf("User not found for the ID %v", id)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("User deleted")
}

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

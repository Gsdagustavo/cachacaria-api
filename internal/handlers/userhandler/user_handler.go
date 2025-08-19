package userhandler

import (
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

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	users := h.UserUseCases.GetAll()

	if len(users) <= 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	prettyJSON, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(prettyJSON)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
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

	log.Printf("Query ID: %v", id)

	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		log.Printf("User not found for the ID %v", id)
		w.Write([]byte(err.Error()))
		return
	}

	user, err := h.UserUseCases.FindById(id)

	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		log.Printf("User not found for the ID %v", id)
		w.Write([]byte(err.Error()))
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
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req models.AddUserRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error adding userhandler: %v", err)
		return
	}

	res, err := h.UserUseCases.Add(req)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error adding userhandler: %v", err)
		return
	}

	//log.Printf("Request: %v", req)

	prettyJSON, err := json.MarshalIndent(req, "", "  ")
	if err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(prettyJSON)

	log.Printf("Response: %v", res)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
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
	w.Write([]byte(`{"message":"userhandler deleted"}`))
	log.Printf("User: %v", id)
}

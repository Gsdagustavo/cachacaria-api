package handlers

import (
	"cachacariaapi/internal/models"
	"cachacariaapi/internal/usecases"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type Handlers struct {
	UserHandler *UserHandler
}

func (h *Handlers) RegisterHandlers(server *http.ServeMux) {
	server.HandleFunc("/users", h.UserHandler.GetUsers)
	server.HandleFunc("/users/id", h.UserHandler.GetUser)
	server.HandleFunc("/users/add", h.UserHandler.AddUser)
}

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

	json.NewEncoder(w).Encode(users)
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

	err = json.NewEncoder(w).Encode(user)

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
		log.Printf("Error adding user: %v", err)
		return
	}

	res, err := h.UserUseCases.Add(req)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error adding user: %v", err)
		return
	}

	log.Printf("Request: %v", req)

	err = json.NewEncoder(w).Encode(res)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error returning response: %v", err)
		return
	}

	log.Printf("Response: %v", res)
}

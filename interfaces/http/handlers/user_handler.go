package handlers

import (
	"cachacariaapi/domain/entities"
	"cachacariaapi/domain/usecases"
	"cachacariaapi/interfaces/http/core"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	UserUseCases *usecases.UserUseCases
}

func NewUserHandler(u *usecases.UserUseCases) *UserHandler {
	return &UserHandler{u}
}

func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) *core.ServerError {
	if apiErr := core.ValidateRequestMethod(r, http.MethodGet); apiErr != nil {
		return apiErr
	}

	users, err := h.UserUseCases.GetAll()

	if err != nil {
		if errors.Is(err, core.ErrNotFound) {
			return &core.ServerError{
				Code:    http.StatusNotFound,
				Message: "no users found",
				Err:     nil,
			}
		}

		return &core.ServerError{
			Code:    http.StatusInternalServerError,
			Message: "could not get users",
			Err:     err,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
	return nil
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) *core.ServerError {
	if apiErr := core.ValidateRequestMethod(r, http.MethodGet); apiErr != nil {
		return apiErr
	}

	vars := mux.Vars(r)
	idStr := vars["id"]
	if idStr == "" {
		return &core.ServerError{
			Code:    http.StatusBadRequest,
			Message: "id is required",
		}
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return &core.ServerError{
			Code:    http.StatusBadRequest,
			Message: "invalid id",
			Err:     err,
		}
	}

	user, err := h.UserUseCases.FindById(id)
	if err != nil {
		if errors.Is(err, core.ErrNotFound) {
			return &core.ServerError{
				Code:    http.StatusNotFound,
				Message: "user not found",
				Err:     err,
			}
		}

		return &core.ServerError{
			Code:    http.StatusInternalServerError,
			Message: core.ErrInternal.Error(),
			Err:     err,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
	return nil
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) *core.ServerError {
	if apiErr := core.ValidateRequestMethod(r, http.MethodDelete); apiErr != nil {
		return apiErr
	}

	vars := mux.Vars(r)
	idStr := vars["id"]
	if idStr == "" {
		return &core.ServerError{
			Code:    http.StatusBadRequest,
			Message: "id is required",
		}
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return &core.ServerError{
			Code:    http.StatusBadRequest,
			Message: "invalid id",
			Err:     err,
		}
	}

	err = h.UserUseCases.Delete(id)
	if err != nil {
		if errors.Is(err, core.ErrNotFound) {
			return &core.ServerError{
				Code:    http.StatusNotFound,
				Message: "user not found",
			}
		}

		return &core.ServerError{
			Code:    http.StatusInternalServerError,
			Message: core.ErrInternal.Error(),
			Err:     err,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	return nil
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) *core.ServerError {
	if apiErr := core.ValidateRequestMethod(r, http.MethodPut); apiErr != nil {
		log.Printf("returning api err %v", apiErr)
		return apiErr
	}

	vars := mux.Vars(r)
	idStr := vars["id"]
	if idStr == "" {
		return &core.ServerError{
			Code:    http.StatusBadRequest,
			Message: "id is required",
		}
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return &core.ServerError{
			Code:    http.StatusBadRequest,
			Message: "invalid id",
			Err:     err,
		}
	}

	var req entities
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return &core.ServerError{
			Code:    http.StatusBadRequest,
			Message: core.ErrBadRequest.Error(),
			Err:     err,
		}
	}

	res, err := h.UserUseCases.Update(req, id)
	if err != nil {
		var apiErr *core.ServerError
		if errors.As(err, &apiErr) {
			return apiErr
		}

		if errors.Is(err, core.ErrNotFound) {
			return &core.ServerError{
				Code:    http.StatusNotFound,
				Message: "user not found",
				Err:     err,
			}
		}

		return &core.ServerError{
			Code:    http.StatusInternalServerError,
			Message: core.ErrInternal.Error(),
			Err:     err,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
	return nil
}

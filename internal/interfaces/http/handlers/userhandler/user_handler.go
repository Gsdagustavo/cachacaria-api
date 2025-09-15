package userhandler

import (
	"cachacariaapi/internal/domain/entities"
	"cachacariaapi/internal/domain/usecases/user"
	"cachacariaapi/internal/interfaces/http/core"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
)

type UserHandler struct {
	UserUseCases *userusecases.UserUseCases
}

func NewUserHandler(u *userusecases.UserUseCases) *UserHandler {
	return &UserHandler{u}
}

func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) *core.ApiError {
	if apiErr := core.ValidateRequestMethod(r, http.MethodGet); apiErr != nil {
		return apiErr
	}

	users, err := h.UserUseCases.GetAll()

	if err != nil {
		if errors.Is(err, core.ErrNotFound) {
			return &core.ApiError{
				Code:    http.StatusNotFound,
				Message: "no users found",
				Err:     nil,
			}
		}

		return &core.ApiError{
			Code:    http.StatusInternalServerError,
			Message: "could not get users",
			Err:     err,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
	return nil
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) *core.ApiError {
	if apiErr := core.ValidateRequestMethod(r, http.MethodGet); apiErr != nil {
		return apiErr
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		return &core.ApiError{
			Code:    http.StatusBadRequest,
			Message: "id is required",
		}
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return &core.ApiError{
			Code:    http.StatusBadRequest,
			Message: "invalid id",
			Err:     err,
		}
	}

	user, err := h.UserUseCases.FindById(id)
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
			Message: core.ErrInternal.Error(),
			Err:     err,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
	return nil
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) *core.ApiError {
	if apiErr := core.ValidateRequestMethod(r, http.MethodDelete); apiErr != nil {
		return apiErr
	}

	query := r.URL.Query()
	if !query.Has("id") {
		return &core.ApiError{
			Code:    http.StatusBadRequest,
			Message: "id is required",
		}
	}

	id, err := strconv.ParseInt(query.Get("id"), 10, 64)
	if err != nil {
		return &core.ApiError{
			Code:    http.StatusBadRequest,
			Message: "invalid id",
		}
	}

	err = h.UserUseCases.Delete(id)
	if err != nil {
		if errors.Is(err, core.ErrNotFound) {
			return &core.ApiError{
				Code:    http.StatusNotFound,
				Message: "user not found",
			}
		}

		return &core.ApiError{
			Code:    http.StatusInternalServerError,
			Message: core.ErrInternal.Error(),
			Err:     err,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entities.UserResponse{ID: id})
	return nil
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) *core.ApiError {
	if apiErr := core.ValidateRequestMethod(r, http.MethodPut); apiErr != nil {
		log.Printf("returning api err %v", apiErr)
		return apiErr
	}

	query := r.URL.Query()

	log.Printf("Query: %v", query)

	hasId := query.Has("id")

	log.Printf("has id: %v", hasId)

	if !hasId {
		return &core.ApiError{
			Code:    http.StatusBadRequest,
			Message: "id is required",
		}
	}

	id, err := strconv.ParseInt(query.Get("id"), 10, 64)

	log.Printf("id: %v", id)
	log.Printf("err: %v", err)

	if err != nil {
		return &core.ApiError{
			Code:    http.StatusBadRequest,
			Message: "invalid id",
		}
	}

	var req entities.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return &core.ApiError{
			Code:    http.StatusBadRequest,
			Message: core.ErrBadRequest.Error(),
			Err:     err,
		}
	}

	res, err := h.UserUseCases.Update(req, id)
	if err != nil {
		var apiErr *core.ApiError
		if errors.As(err, &apiErr) {
			return apiErr
		}

		if errors.Is(err, core.ErrNotFound) {
			return &core.ApiError{
				Code:    http.StatusNotFound,
				Message: "user not found",
				Err:     err,
			}
		}

		return &core.ApiError{
			Code:    http.StatusInternalServerError,
			Message: core.ErrInternal.Error(),
			Err:     err,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
	return nil
}

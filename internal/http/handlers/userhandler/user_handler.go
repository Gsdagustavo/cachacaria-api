package userhandler

import (
	"cachacariaapi/internal/http/core"
	"cachacariaapi/internal/models"
	"cachacariaapi/internal/usecases"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
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

func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) *core.ApiError {
	if apiErr := ValidateRequestMethod(r, http.MethodGet); apiErr != nil {
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
	if apiErr := ValidateRequestMethod(r, http.MethodGet); apiErr != nil {
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

func (h *UserHandler) AddUser(w http.ResponseWriter, r *http.Request) *core.ApiError {
	if apiErr := ValidateRequestMethod(r, http.MethodPost); apiErr != nil {
		return apiErr
	}

	var req models.UserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return &core.ApiError{
			Code:    http.StatusBadRequest,
			Message: core.ErrBadRequest.Error(),
			Err:     err,
		}
	}

	res, err := h.UserUseCases.Add(req)
	if err != nil {
		var apiErr *core.ApiError
		if errors.As(err, &apiErr) {
			return apiErr
		}

		if errors.Is(err, core.ErrBadRequest) {
			return &core.ApiError{
				Code:    http.StatusBadRequest,
				Message: core.ErrBadRequest.Error(),
				Err:     err,
			}
		}

		if errors.Is(err, core.ErrConflict) {
			return &core.ApiError{
				Code:    http.StatusConflict,
				Message: core.ErrConflict.Error(),
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

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) *core.ApiError {
	if apiErr := ValidateRequestMethod(r, http.MethodDelete); apiErr != nil {
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
	json.NewEncoder(w).Encode(models.UserResponse{ID: id})
	return nil
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) *core.ApiError {
	if apiErr := ValidateRequestMethod(r, http.MethodDelete); apiErr != nil {
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

	var req models.UserRequest
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

package producthandler

import (
	"cachacariaapi/internal/domain/entities"
	"cachacariaapi/internal/domain/usecases/product"
	"cachacariaapi/internal/interfaces/http/core"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
)

type ProductHandler struct {
	ProductUseCases *product.ProductUseCases
}

func NewProductHandler(productUseCases *product.ProductUseCases) *ProductHandler {
	return &ProductHandler{productUseCases}
}

func (h *ProductHandler) Add(w http.ResponseWriter, r *http.Request) *core.ApiError {
	if apiErr := core.ValidateRequestMethod(r, http.MethodPost); apiErr != nil {
		return apiErr.WithError("response handler / add response")
	}

	if err := r.ParseMultipartForm(10 >> 20); err != nil {
		return (&core.ApiError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Err:     err,
		}).WithError("product handler / parse multipart form failed: " + err.Error())
	}

	name := r.FormValue("name")
	description := r.FormValue("description")
	price, _ := strconv.ParseFloat(r.FormValue("price"), 64)
	stock, _ := strconv.Atoi(r.FormValue("stock"))

	request := entities.AddProductRequest{
		Name:        name,
		Description: description,
		Price:       float32(price),
		Stock:       stock,
	}

	photos := r.MultipartForm.File["photos"]

	response, err := h.ProductUseCases.Add(request, photos)
	if err != nil {
		log.Printf("error adding response: %v", err)

		if errors.Is(err, core.ErrConflict) {
			return (&core.ApiError{
				Code:    http.StatusConflict,
				Err:     err,
				Message: "response already exists",
			}).WithError("response handler / add response")
		}

		if errors.Is(err, core.ErrBadRequest) {
			return (&core.ApiError{
				Code:    http.StatusBadRequest,
				Err:     err,
				Message: "invalid response data",
			}).WithError("response handler / add response")
		}

		return (&core.ApiError{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
			Err:     err,
		}).WithError("product handler / add response")
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	return nil
}

func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) *core.ApiError {
	if apiErr := core.ValidateRequestMethod(r, http.MethodGet); apiErr != nil {
		return apiErr.WithError("product handler / get all")
	}

	users, err := h.ProductUseCases.GetAll()

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

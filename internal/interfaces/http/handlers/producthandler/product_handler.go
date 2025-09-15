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

	"github.com/gorilla/mux"
)

type ProductHandler struct {
	ProductUseCases *product.ProductUseCases
}

func NewProductHandler(productUseCases *product.ProductUseCases) *ProductHandler {
	return &ProductHandler{productUseCases}
}

func (h *ProductHandler) Add(w http.ResponseWriter, r *http.Request) *core.ApiError {
	if apiErr := core.ValidateRequestMethod(r, http.MethodPost); apiErr != nil {
		return apiErr.WithError("products handler / add product")
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

	response, err := h.ProductUseCases.AddProduct(request, photos)
	if err != nil {
		log.Printf("error adding products: %v", err)

		if errors.Is(err, core.ErrConflict) {
			return (&core.ApiError{
				Code:    http.StatusConflict,
				Err:     err,
				Message: "product already exists",
			}).WithError("products handler / add products")
		}

		if errors.Is(err, core.ErrBadRequest) {
			return (&core.ApiError{
				Code:    http.StatusBadRequest,
				Err:     err,
				Message: "invalid response data",
			}).WithError("products handler / add products")
		}

		return (&core.ApiError{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
			Err:     err,
		}).WithError("product handler / add products")
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	return nil
}

func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) *core.ApiError {
	if apiErr := core.ValidateRequestMethod(r, http.MethodGet); apiErr != nil {
		return apiErr.WithError("product handler / get all")
	}

	products, err := h.ProductUseCases.GetAll()

	if err != nil {
		if errors.Is(err, core.ErrNotFound) {
			return (&core.ApiError{
				Code:    http.StatusNotFound,
				Message: "no products found",
				Err:     nil,
			}).WithError("product handler / get all")
		}

		return (&core.ApiError{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
			Err:     err,
		}).WithError("product handler / get products")
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
	return nil
}

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) *core.ApiError {
	if apiErr := core.ValidateRequestMethod(r, http.MethodGet); apiErr != nil {
		return apiErr.WithError("product handler / get all")
	}

	vars := mux.Vars(r)
	idStr := vars["id"]
	if idStr == "" {
		return (&core.ApiError{
			Code:    http.StatusBadRequest,
			Message: "id is required",
		}).WithError("product handler / get product")
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return (&core.ApiError{
			Code:    http.StatusBadRequest,
			Message: "invalid id",
			Err:     err,
		}).WithError("product handler / get product")
	}

	product, err := h.ProductUseCases.GetProduct(id)
	if err != nil {
		if errors.Is(err, core.ErrNotFound) {
			return (&core.ApiError{
				Code:    http.StatusNotFound,
				Message: "product not found",
				Err:     nil,
			}).WithError("product handler / get product")
		}

		return (&core.ApiError{
			Code:    http.StatusInternalServerError,
			Message: "could not get product",
			Err:     err,
		}).WithError("product handler / get product")
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
	return nil
}

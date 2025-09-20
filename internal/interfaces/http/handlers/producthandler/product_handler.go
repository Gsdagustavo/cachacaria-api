package producthandler

import (
	"cachacariaapi/internal/domain/entities"
	"cachacariaapi/internal/domain/usecases/product"
	"cachacariaapi/internal/interfaces/http/core"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

const maxImagesMemory = 20 << 20

type ProductHandler struct {
	ProductUseCases *product.ProductUseCases
}

func NewProductHandler(productUseCases *product.ProductUseCases) *ProductHandler {
	return &ProductHandler{productUseCases}
}

func (h *ProductHandler) Add(w http.ResponseWriter, r *http.Request) *core.ServerError {
	if apiErr := core.ValidateRequestMethod(r, http.MethodPost); apiErr != nil {
		return apiErr.WithError("products handler / add product")
	}

	if err := r.ParseMultipartForm(maxImagesMemory); err != nil {
		return (&core.ServerError{
			Code:    http.StatusBadRequest,
			Message: "product images form exceeds maximum memory limit",
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
		log.Printf("error adding product: %v", err)

		if errors.Is(err, core.ErrConflict) {
			return (&core.ServerError{
				Code:    http.StatusConflict,
				Err:     err,
				Message: "product already exists",
			}).WithError("products handler / add product")
		}

		if errors.Is(err, core.ErrBadRequest) {
			return (&core.ServerError{
				Code:    http.StatusBadRequest,
				Err:     err,
				Message: "bad request",
			}).WithError("products handler / add product")
		}

		return (&core.ServerError{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
			Err:     err,
		}).WithError("product handler / add product")
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	return nil
}

func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) *core.ServerError {
	if apiErr := core.ValidateRequestMethod(r, http.MethodGet); apiErr != nil {
		return apiErr.WithError("product handler / get all")
	}

	products, err := h.ProductUseCases.GetAll()

	if err != nil {
		if errors.Is(err, core.ErrNotFound) {
			return (&core.ServerError{
				Code:    http.StatusNotFound,
				Message: "no products found",
				Err:     nil,
			}).WithError("product handler / get all")
		}

		return (&core.ServerError{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
			Err:     err,
		}).WithError("product handler / get products")
	}

	baseURL := os.Getenv("BASE_URL")
	for i := range products {
		var fullURLs []string
		for _, filename := range products[i].Photos {
			fullURLs = append(fullURLs, fmt.Sprintf("%s/images/%s", baseURL, filename))
		}
		products[i].Photos = fullURLs
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
	return nil
}

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) *core.ServerError {
	if apiErr := core.ValidateRequestMethod(r, http.MethodGet); apiErr != nil {
		return apiErr.WithError("prod handler / get all")
	}

	vars := mux.Vars(r)
	idStr := vars["id"]
	if idStr == "" {
		return (&core.ServerError{
			Code:    http.StatusBadRequest,
			Message: "id is required",
		}).WithError("prod handler / get prod")
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return (&core.ServerError{
			Code:    http.StatusBadRequest,
			Message: "invalid id",
			Err:     err,
		}).WithError("prod handler / get prod")
	}

	prod, err := h.ProductUseCases.GetProduct(id)
	if err != nil {
		if errors.Is(err, core.ErrNotFound) {
			return (&core.ServerError{
				Code:    http.StatusNotFound,
				Message: "prod not found",
				Err:     nil,
			}).WithError("prod handler / get prod")
		}

		return (&core.ServerError{
			Code:    http.StatusInternalServerError,
			Message: "could not get prod",
			Err:     err,
		}).WithError("prod handler / get prod")
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(prod)
	return nil
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) *core.ServerError {
	if apiErr := core.ValidateRequestMethod(r, http.MethodDelete); apiErr != nil {
		return apiErr.WithError("prod handler / delete")
	}

	vars := mux.Vars(r)
	idStr := vars["id"]
	if idStr == "" {
		return (&core.ServerError{
			Code:    http.StatusBadRequest,
			Message: "id is required",
		}).WithError("prod handler / delete")
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return (&core.ServerError{
			Code:    http.StatusBadRequest,
			Message: "invalid id",
			Err:     err,
		}).WithError("prod handler / delete")
	}

	res, err := h.ProductUseCases.DeleteProduct(id)
	if err != nil {
		if errors.Is(err, core.ErrNotFound) {
			return (&core.ServerError{
				Code:    http.StatusNotFound,
				Err:     nil,
				Message: "product not found",
			}).WithError("prod handler / delete")
		}

		return (&core.ServerError{
			Code:    http.StatusInternalServerError,
			Message: "could not delete product",
			Err:     err,
		}).WithError("prod handler / delete")
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
	return nil
}

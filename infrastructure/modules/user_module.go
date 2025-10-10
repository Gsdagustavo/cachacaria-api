package modules

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

const maxImagesMemory = 20 << 20

const defaultPagePagination = 1
const defaultLimitPagination = 20

type ProductModule struct {
	ProductUseCases *usecases.ProductUseCases
	name            string
	path            string
}

func NewProductModule(productUseCases *usecases.ProductUseCases) *ProductModule {
	return &ProductModule{
		ProductUseCases: productUseCases,
		name:            "auth",
		path:            "/auth",
	}
}

func (m ProductModule) Name() string {
	return m.name
}

func (m ProductModule) Path() string {
	return m.path
}

func (m ProductModule) add(w http.ResponseWriter, r *http.Request) *core.ServerError {
	if apiErr := core.ValidateRequestMethod(r, http.MethodPost); apiErr != nil {
		return apiErr.WithError("products handler / add product")
	}

	if err := r.ParseMultipartForm(maxImagesMemory); err != nil {
		return (&core.ServerError{
			Code:    http.StatusBadRequest,
			Message: "Imagens excedem o máximo de memória permitido",
			Err:     err,
		}).WithError("product handler / parse multipart form failed: " + err.Error())
	}

	name := r.FormValue("name")
	description := r.FormValue("description")
	price, _ := strconv.ParseFloat(r.FormValue("price"), 64)
	stock, _ := strconv.Atoi(r.FormValue("stock"))
	photos := r.MultipartForm.File["photos"]

	request := entities.AddProductRequest{
		Name:        name,
		Description: description,
		Price:       float32(price),
		Stock:       stock,
		Photos:      photos,
	}

	response, err := m.ProductUseCases.AddProduct(request)
	if err != nil {
		log.Printf("error adding product: %v", err)

		if errors.Is(err, core.ErrConflict) {
			return (&core.ServerError{
				Code:    http.StatusConflict,
				Message: "Este produto já existew",
			}).WithError("products handler / add product")
		}

		if errors.Is(err, core.ErrBadRequest) {
			return (&core.ServerError{
				Code:    http.StatusBadRequest,
				Message: "Requisição inválida",
			}).WithError("products handler / add product")
		}

		return (&core.ServerError{
			Code:    http.StatusInternalServerError,
			Message: "Erro interno no servidor",
		}).WithError("product handler / add product")
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	return nil
}

func (m ProductModule) getAll(w http.ResponseWriter, r *http.Request) *core.ServerError {
	if apiErr := core.ValidateRequestMethod(r, http.MethodGet); apiErr != nil {
		return apiErr.WithError("product handler / get all")
	}

	var limit, page int

	query := r.URL.Query()
	limitStr := query.Get("limit")
	pageStr := query.Get("page")

	if limitStr == "" {
		limit = defaultLimitPagination
	} else {
		parsed, err := strconv.ParseInt(limitStr, 10, 32)
		if err != nil {
			return (&core.ServerError{
				Code:    http.StatusBadRequest,
				Message: "Limite inválido",
			}).WithError("prod handler / get prod paginated")
		}

		limit = int(parsed)
	}

	if pageStr == "" {
		page = defaultPagePagination
	} else {
		parsed, err := strconv.ParseInt(pageStr, 10, 32)
		if err != nil {
			return (&core.ServerError{
				Code:    http.StatusBadRequest,
				Message: "Página inválida",
			}).WithError("prod handler / get prod paginated")
		}

		page = int(parsed)
	}

	products, err := m.ProductUseCases.GetAll(limit, page)

	if err != nil {
		if errors.Is(err, core.ErrNotFound) {
			return (&core.ServerError{
				Code:    http.StatusNotFound,
				Message: "Nenhum produto encontrado",
				Err:     nil,
			}).WithError("product handler / get all")
		}

		return (&core.ServerError{
			Code:    http.StatusInternalServerError,
			Message: "Erro interno no servidor",
		}).WithError("product handler / get products")
	}

	//baseURL := os.Getenv("BASE_URL")
	//for i := range products {
	//	var fullURLs []string
	//	for _, filename := range products[i].Photos {
	//		fullURLs = append(fullURLs, fmt.Sprintf("%s/images/%s", baseURL, filename))
	//	}
	//	products[i].Photos = fullURLs
	//}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
	return nil
}

func (m ProductModule) get(w http.ResponseWriter, r *http.Request) *core.ServerError {
	if apiErr := core.ValidateRequestMethod(r, http.MethodGet); apiErr != nil {
		return apiErr.WithError("prod handler / get all")
	}

	vars := mux.Vars(r)
	idStr := vars["id"]
	if idStr == "" {
		return (&core.ServerError{
			Code:    http.StatusBadRequest,
			Message: "ID do produto inválido",
		}).WithError("prod handler / get prod")
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return (&core.ServerError{
			Code:    http.StatusBadRequest,
			Message: "ID do produto inválido",
		}).WithError("prod handler / get prod")
	}

	prod, err := m.ProductUseCases.GetProduct(id)
	if err != nil {
		if errors.Is(err, core.ErrNotFound) {
			return (&core.ServerError{
				Code:    http.StatusNotFound,
				Message: "Produto não encontrado",
			}).WithError("prod handler / get prod")
		}

		return (&core.ServerError{
			Code:    http.StatusInternalServerError,
			Message: "Erro interno no servidor",
		}).WithError("prod handler / get prod")
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(prod)
	return nil
}

func (m ProductModule) delete(w http.ResponseWriter, r *http.Request) *core.ServerError {
	if apiErr := core.ValidateRequestMethod(r, http.MethodDelete); apiErr != nil {
		return apiErr.WithError("prod handler / delete")
	}

	vars := mux.Vars(r)
	idStr := vars["id"]
	if idStr == "" {
		return (&core.ServerError{
			Code:    http.StatusBadRequest,
			Message: "ID do produto inválido",
		}).WithError("prod handler / delete")
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return (&core.ServerError{
			Code:    http.StatusBadRequest,
			Message: "ID do produto inválido",
		}).WithError("prod handler / delete")
	}

	res, err := m.ProductUseCases.DeleteProduct(id)
	if err != nil {
		if errors.Is(err, core.ErrNotFound) {
			return (&core.ServerError{
				Code:    http.StatusNotFound,
				Message: "Produto não encontrado",
			}).WithError("prod handler / delete")
		}

		return (&core.ServerError{
			Code:    http.StatusInternalServerError,
			Message: "Erro interno no servidor",
		}).WithError("prod handler / delete")
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
	return nil
}

func (m ProductModule) update(w http.ResponseWriter, r *http.Request) *core.ServerError {
	if apiErr := core.ValidateRequestMethod(r, http.MethodPatch); apiErr != nil {
		return apiErr.WithError("product handler / update product")
	}

	vars := mux.Vars(r)
	idStr := vars["id"]
	if idStr == "" {
		return (&core.ServerError{
			Code:    http.StatusBadRequest,
			Message: "ID do produto inválido",
		}).WithError("prod handler / update")
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return (&core.ServerError{
			Code:    http.StatusBadRequest,
			Message: "ID do produto inválido",
			Err:     err,
		}).WithError("prod handler / update")
	}

	if err := r.ParseMultipartForm(maxImagesMemory); err != nil {
		return (&core.ServerError{
			Code:    http.StatusBadRequest,
			Message: "",
			Err:     err,
		}).WithError("product handler / parse multipart form failed: " + err.Error())
	}

	name := r.FormValue("name")
	description := r.FormValue("description")
	price, _ := strconv.ParseFloat(r.FormValue("price"), 64)
	stock, _ := strconv.Atoi(r.FormValue("stock"))

	finalPrice := float32(price)

	// photos := r.MultipartForm.File["photos"]

	request := entities.UpdateProductRequest{
		Name:        &name,
		Description: &description,
		Price:       &finalPrice,
		Stock:       &stock,
	}

	if err = json.NewDecoder(r.Body).Decode(&request); err != nil {
		return (&core.ServerError{
			Code: http.StatusBadRequest,
		}).WithError("prod handler / update")
	}

	res, err := m.ProductUseCases.UpdateProduct(id, request)
	if err != nil {
		if errors.Is(err, core.ErrNotFound) {
			return (&core.ServerError{
				Code:    http.StatusNotFound,
				Message: "Produto não encontrado",
			}).WithError("prod handler / update")
		}

		return (&core.ServerError{
			Code:    http.StatusInternalServerError,
			Message: "Erro interno no servidor",
		}).WithError("prod handler / update")
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
	return nil
}

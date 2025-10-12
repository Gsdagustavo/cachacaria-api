package modules

import (
	"cachacariaapi/domain/entities"
	"cachacariaapi/domain/usecases"
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
		name:            "product",
		path:            "/product",
	}
}

func (m ProductModule) Name() string {
	return m.name
}

func (m ProductModule) Path() string {
	return m.path
}

func (m ProductModule) RegisterRoutes(router *mux.Router) {
	routes := []ModuleRoute{
		{
			Name:    "Add",
			Path:    "",
			Handler: m.add,
			Methods: []string{http.MethodPost},
		},
		{
			Name:    "GetAll",
			Path:    "",
			Handler: m.getAll,
			Methods: []string{http.MethodGet},
		},
		{
			Name:    "Get",
			Path:    "/{id}",
			Handler: m.get,
			Methods: []string{http.MethodGet},
		},
		{
			Name:    "Update",
			Path:    "/{id}",
			Handler: m.update,
			Methods: []string{http.MethodPatch},
		},
		{
			Name:    "Delete",
			Path:    "/{id}",
			Handler: m.delete,
			Methods: []string{http.MethodDelete},
		},
	}

	for _, route := range routes {
		router.HandleFunc(m.path+route.Path, route.Handler).Methods(route.Methods...)
	}
}

func (m ProductModule) add(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(maxImagesMemory); err != nil {
		res := entities.ServerResponse{
			Code:    http.StatusBadRequest,
			Message: "Imagens excedem o máximo de memória permitido",
		}
		res.WriteHTTP(w)
		return
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

		var res entities.ServerResponse
		if errors.Is(err, entities.ErrConflict) {
			res = entities.ServerResponse{
				Code:    http.StatusConflict,
				Message: "Este produto já existe",
			}
		} else if errors.Is(err, entities.ErrBadRequest) {
			res = entities.ServerResponse{
				Code:    http.StatusBadRequest,
				Message: "Requisição inválida",
			}
		} else {
			res = entities.ServerResponse{
				Code:    http.StatusInternalServerError,
				Message: "Erro interno no servidor",
			}
		}
		res.WriteHTTP(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (m ProductModule) getAll(w http.ResponseWriter, r *http.Request) {
	var limit, page int

	query := r.URL.Query()
	limitStr := query.Get("limit")
	pageStr := query.Get("page")

	if limitStr == "" {
		limit = defaultLimitPagination
	} else {
		parsed, err := strconv.ParseInt(limitStr, 10, 32)
		if err != nil {
			res := entities.ServerResponse{
				Code:    http.StatusBadRequest,
				Message: "Limite inválido",
			}
			res.WriteHTTP(w)
			return
		}
		limit = int(parsed)
	}

	if pageStr == "" {
		page = defaultPagePagination
	} else {
		parsed, err := strconv.ParseInt(pageStr, 10, 32)
		if err != nil {
			res := entities.ServerResponse{
				Code:    http.StatusBadRequest,
				Message: "Página inválida",
			}
			res.WriteHTTP(w)
			return
		}
		page = int(parsed)
	}

	products, err := m.ProductUseCases.GetAll(limit, page)
	if err != nil {
		var res entities.ServerResponse
		if errors.Is(err, entities.ErrNotFound) {
			res = entities.ServerResponse{
				Code:    http.StatusNotFound,
				Message: "Nenhum produto encontrado",
			}
		} else {
			res = entities.ServerResponse{
				Code:    http.StatusInternalServerError,
				Message: "Erro interno no servidor",
			}
		}
		res.WriteHTTP(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func (m ProductModule) get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	var res entities.ServerResponse
	if idStr == "" {
		res = entities.ServerResponse{
			Code:    http.StatusBadRequest,
			Message: "ID do produto inválido",
		}
		res.WriteHTTP(w)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		res = entities.ServerResponse{
			Code:    http.StatusBadRequest,
			Message: "ID do produto inválido",
		}
		res.WriteHTTP(w)
		return
	}

	prod, err := m.ProductUseCases.GetProduct(id)
	if err != nil {
		if errors.Is(err, entities.ErrNotFound) {
			res = entities.ServerResponse{
				Code:    http.StatusNotFound,
				Message: "Produto não encontrado",
			}
		} else {
			res = entities.ServerResponse{
				Code:    http.StatusInternalServerError,
				Message: "Erro interno no servidor",
			}
		}
		res.WriteHTTP(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(prod)
}

func (m ProductModule) delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	var res entities.ServerResponse
	if idStr == "" {
		res = entities.ServerResponse{
			Code:    http.StatusBadRequest,
			Message: "ID do produto inválido",
		}
		res.WriteHTTP(w)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		res = entities.ServerResponse{
			Code:    http.StatusBadRequest,
			Message: "ID do produto inválido",
		}
		res.WriteHTTP(w)
		return
	}

	delRes, err := m.ProductUseCases.DeleteProduct(id)
	if err != nil {
		if errors.Is(err, entities.ErrNotFound) {
			res = entities.ServerResponse{
				Code:    http.StatusNotFound,
				Message: "Produto não encontrado",
			}
		} else {
			res = entities.ServerResponse{
				Code:    http.StatusInternalServerError,
				Message: "Erro interno no servidor",
			}
		}
		res.WriteHTTP(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(delRes)
}

func (m ProductModule) update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	var res entities.ServerResponse
	if idStr == "" {
		res = entities.ServerResponse{
			Code:    http.StatusBadRequest,
			Message: "ID do produto inválido",
		}
		res.WriteHTTP(w)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		res = entities.ServerResponse{
			Code:    http.StatusBadRequest,
			Message: "ID do produto inválido",
		}
		res.WriteHTTP(w)
		return
	}

	var request entities.UpdateProductRequest
	if err = json.NewDecoder(r.Body).Decode(&request); err != nil {
		res = entities.ServerResponse{
			Code:    http.StatusBadRequest,
			Message: "Requisição inválida. Certifique-se de usar application/json.",
		}
		res.WriteHTTP(w)
		return
	}

	updateRes, err := m.ProductUseCases.UpdateProduct(id, request)
	if err != nil {
		if errors.Is(err, entities.ErrNotFound) {
			res = entities.ServerResponse{
				Code:    http.StatusNotFound,
				Message: "Produto não encontrado",
			}
		} else {
			res = entities.ServerResponse{
				Code:    http.StatusInternalServerError,
				Message: "Erro interno no servidor",
			}
		}
		res.WriteHTTP(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updateRes)
}

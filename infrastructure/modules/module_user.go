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

type UserModule struct {
	userUseCases *usecases.UserUseCases
	name         string
	path         string
}

func NewUserModule(userUseCases *usecases.UserUseCases) *UserModule {
	return &UserModule{
		userUseCases: userUseCases,
		name:         "user",
		path:         "/user",
	}
}

func (m UserModule) Name() string {
	return m.name
}

func (m UserModule) Path() string {
	return m.path
}

func (m UserModule) RegisterRoutes(router *mux.Router) {
	routes := []ModuleRoute{
		{
			Name:    "GetAll",
			Path:    "",
			Handler: m.GetAll,
			Methods: []string{http.MethodGet},
		},
		{
			Name:    "GetUser",
			Path:    "/{id}",
			Handler: m.GetUser,
			Methods: []string{http.MethodGet},
		},
		{
			Name:    "UpdateUser",
			Path:    "/{id}",
			Handler: m.UpdateUser,
			Methods: []string{http.MethodPut},
		},
		{
			Name:    "DeleteUser",
			Path:    "/{id}",
			Handler: m.DeleteUser,
			Methods: []string{http.MethodDelete},
		},
	}

	for _, route := range routes {
		router.HandleFunc(m.path+route.Path, route.Handler).Methods(route.Methods...)
	}
}

func (m UserModule) GetAll(w http.ResponseWriter, r *http.Request) {
	users, err := m.userUseCases.GetAll()

	if err != nil {
		var res entities.ServerResponse
		if errors.Is(err, entities.ErrNotFound) || errors.Is(err, entities.ErrUserNotFound) {
			res = entities.ServerResponse{
				Code:    http.StatusNotFound,
				Message: "Nenhum usuário encontrado",
			}
		} else {
			log.Printf("error getting all users: %v", err)
			res = entities.ServerResponse{
				Code:    http.StatusInternalServerError,
				Message: entities.ErrInternal.Error(),
			}
		}
		res.WriteHTTP(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (m UserModule) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	var res entities.ServerResponse
	if idStr == "" {
		res = entities.ServerResponse{
			Code:    http.StatusBadRequest,
			Message: "ID de usuário inválido",
		}
		res.WriteHTTP(w)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		res = entities.ServerResponse{
			Code:    http.StatusBadRequest,
			Message: "ID de usuário inválido",
		}
		res.WriteHTTP(w)
		return
	}

	user, err := m.userUseCases.FindById(id)
	if err != nil {
		if errors.Is(err, entities.ErrNotFound) || errors.Is(err, entities.ErrUserNotFound) {
			res = entities.ServerResponse{
				Code:    http.StatusNotFound,
				Message: "Usuário não encontrado",
			}
		} else {
			log.Printf("error finding user %d: %v", id, err)
			res = entities.ServerResponse{
				Code:    http.StatusInternalServerError,
				Message: entities.ErrInternal.Error(),
			}
		}
		res.WriteHTTP(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (m UserModule) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	var res entities.ServerResponse
	if idStr == "" {
		res = entities.ServerResponse{
			Code:    http.StatusBadRequest,
			Message: "ID de usuário inválido",
		}
		res.WriteHTTP(w)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		res = entities.ServerResponse{
			Code:    http.StatusBadRequest,
			Message: "ID de usuário inválido",
		}
		res.WriteHTTP(w)
		return
	}

	err = m.userUseCases.Delete(id)
	if err != nil {
		if errors.Is(err, entities.ErrNotFound) || errors.Is(err, entities.ErrUserNotFound) {
			res = entities.ServerResponse{
				Code:    http.StatusNotFound,
				Message: "Usuário não encontrado",
			}
		} else {
			log.Printf("error deleting user %d: %v", id, err)
			res = entities.ServerResponse{
				Code:    http.StatusInternalServerError,
				Message: entities.ErrInternal.Error(),
			}
		}
		res.WriteHTTP(w)
		return
	}

	response := struct {
		Message string `json:"message"`
		Status  int    `json:"status"`
	}{
		Message: "Usuário excluído com sucesso",
		Status:  http.StatusOK,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (m UserModule) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	var res entities.ServerResponse
	if idStr == "" {
		res = entities.ServerResponse{
			Code:    http.StatusBadRequest,
			Message: "ID de usuário inválido",
		}
		res.WriteHTTP(w)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		res = entities.ServerResponse{
			Code:    http.StatusBadRequest,
			Message: "ID de usuário inválido",
		}
		res.WriteHTTP(w)
		return
	}

	var req entities.User
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		res = entities.ServerResponse{
			Code:    http.StatusBadRequest,
			Message: "Requisição inválida. Certifique-se de usar application/json.",
		}
		res.WriteHTTP(w)
		return
	}

	err = m.userUseCases.Update(req, id)
	if err != nil {
		var apiErr *entities.ServerError
		if errors.As(err, &apiErr) {
			apiErr.WriteHTTP(w)
			return
		}

		if errors.Is(err, entities.ErrNotFound) || errors.Is(err, entities.ErrUserNotFound) {
			res = entities.ServerResponse{
				Code:    http.StatusNotFound,
				Message: "Usuário não encontrado",
			}
		} else {
			log.Printf("error updating user %d: %v", id, err)
			res = entities.ServerResponse{
				Code:    http.StatusInternalServerError,
				Message: entities.ErrInternal.Error(),
			}
		}
		res.WriteHTTP(w)
		return
	}

	response := struct {
		Message string `json:"message"`
		Status  int    `json:"status"`
	}{
		Message: "Usuário atualizado com sucesso",
		Status:  http.StatusOK,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

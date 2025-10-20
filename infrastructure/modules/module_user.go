package modules

import (
	"cachacariaapi/domain/entities"
	"cachacariaapi/domain/usecases"
	"cachacariaapi/infrastructure/util"
	"log/slog"
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

func (m UserModule) GetAll(w http.ResponseWriter, _ *http.Request) {
	users, err := m.userUseCases.GetAll()

	if err != nil {
		slog.Error("failed to get all users", "cause", err)
		util.WriteInternalError(w)
		return
	}

	util.Write(w, users)
}

func (m UserModule) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	if idStr == "" {
		util.WriteBadRequest(w)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		util.WriteBadRequest(w)
		return
	}

	user, err := m.userUseCases.FindById(id)
	if err != nil {
		slog.Error("failed to get user by id", "cause", err)
		util.WriteInternalError(w)
	}

	if user == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	util.Write(w, user)
}

func (m UserModule) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	if idStr == "" {
		util.WriteBadRequest(w)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		util.WriteBadRequest(w)
		return
	}

	err = m.userUseCases.Delete(id)
	if err != nil {
		slog.Error("failed to delete user", "cause", err)
		util.WriteInternalError(w)
	}

	w.WriteHeader(http.StatusOK)
}

func (m UserModule) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		util.WriteBadRequest(w)
		return
	}

	var req entities.User
	err = util.Read(r, &req)
	if err != nil {
		util.WriteBadRequest(w)
		return
	}

	err = m.userUseCases.Update(req, id)
	if err != nil {
		slog.Error("failed to update user", "cause", err)
		util.WriteInternalError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
}

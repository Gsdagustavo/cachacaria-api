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

type moduleUser struct {
	userUseCases usecases.UserUseCases
	authUseCases usecases.AuthUseCases
	name         string
	path         string
}

func NewUserModule(userUseCases usecases.UserUseCases, authUseCases usecases.AuthUseCases) Module {
	return moduleUser{
		userUseCases: userUseCases,
		authUseCases: authUseCases,
		name:         "user",
		path:         "/user",
	}
}

func (m moduleUser) Name() string {
	return m.name
}

func (m moduleUser) Path() string {
	return m.path
}

func (m moduleUser) RegisterRoutes(router *mux.Router) {
	routes := []ModuleRoute{
		{
			Name:    "GetAll",
			Path:    "",
			Handler: m.GetAll,
			Methods: []string{http.MethodGet},
		},
		{
			Name:    "UpdateUser",
			Path:    "/",
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

func (m moduleUser) GetAll(w http.ResponseWriter, _ *http.Request) {
	users, err := m.userUseCases.GetAll()

	if err != nil {
		slog.Error("failed to get all users", "cause", err)
		util.WriteInternalError(w)
		return
	}

	util.Write(w, users)
}

func (m moduleUser) DeleteUser(w http.ResponseWriter, r *http.Request) {
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

func (m moduleUser) UpdateUser(w http.ResponseWriter, r *http.Request) {
	token := util.GetAuthTokenFromRequest(r)

	user, err := m.authUseCases.GetUserByAuthToken(token)
	if err != nil {
		slog.Error("failed to get user by auth token", "cause", err)
		util.WriteInternalError(w)
		return
	}

	if user == nil {
		util.WriteUnauthorized(w)
		return
	}

	var req entities.User
	err = util.Read(r, &req)
	if err != nil {
		util.WriteBadRequest(w)
		return
	}

	status, err := m.userUseCases.Update(r.Context(), req, int64(user.ID))
	if err != nil {
		slog.Error("failed to update user", "cause", err)
		util.WriteInternalError(w)
		return
	}

	res := util.ServerResponse{
		Status:  status.Int(),
		Message: status.String(),
	}
	util.Write(w, res)
}

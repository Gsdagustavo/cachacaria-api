package modules

import (
	"cachacariaapi/infrastructure/util"
	"net/http"

	"github.com/gorilla/mux"
)

type HealthModule struct {
	name string
	path string
}

func NewHealthModule() Module {
	return HealthModule{
		name: "health",
		path: "/health",
	}
}

func (h HealthModule) Name() string {
	return h.name
}

func (h HealthModule) Path() string {
	return h.path
}

func (h HealthModule) RegisterRoutes(router *mux.Router) {
	routes := []ModuleRoute{
		{
			Name:    "health",
			Path:    h.path,
			Handler: h.health,
			Methods: []string{http.MethodGet},
		},
	}

	for _, route := range routes {
		router.HandleFunc(route.Path, route.Handler).Methods(route.Methods...)
	}
}

func (h HealthModule) health(w http.ResponseWriter, _ *http.Request) {
	response := util.ServerResponse{
		Status:  http.StatusOK,
		Message: "Server is healthy",
	}

	util.Write(w, response)
}

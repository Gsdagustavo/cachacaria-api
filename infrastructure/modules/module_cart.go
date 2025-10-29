package modules

import (
	"cachacariaapi/domain/usecases"

	"github.com/gorilla/mux"
)

type CartModule struct {
	CartUseCases *usecases.CartUseCases
	name         string
	path         string
}

func NewCartModule(cartUseCases *usecases.CartUseCases) *CartModule {
	return &CartModule{
		CartUseCases: cartUseCases,
		name:         "cart",
		path:         "/cart",
	}
}

func (m CartModule) Name() string {
	return m.name
}

func (m CartModule) Path() string {
	return m.path
}

func (m CartModule) RegisterRoutes(router *mux.Router) {
	//routes := []ModuleRoute{
	//	{
	//		Name:    "Add",
	//		Path:    "",
	//		Handler: m.add,
	//		Methods: []string{http.MethodPost},
	//	},
	//	{
	//		Name:    "GetAll",
	//		Path:    "",
	//		Handler: m.getAll,
	//		Methods: []string{http.MethodGet},
	//	},
	//	{
	//		Name:    "Get",
	//		Path:    "/{id}",
	//		Handler: m.get,
	//		Methods: []string{http.MethodGet},
	//	},
	//	{
	//		Name:    "Update",
	//		Path:    "/{id}",
	//		Handler: m.update,
	//		Methods: []string{http.MethodPatch},
	//	},
	//	{
	//		Name:    "Delete",
	//		Path:    "/{id}",
	//		Handler: m.delete,
	//		Methods: []string{http.MethodDelete},
	//	},
	//}

	//for _, route := range routes {
	//	router.HandleFunc(m.path+route.Path, route.Handler).Methods(route.Methods...)
	//}
}

package modules

import (
	"cachacariaapi/domain/usecases"
	"cachacariaapi/infrastructure/middleware"
	"cachacariaapi/infrastructure/util"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
)

type OrderModule struct {
	cartUseCases usecases.CartUseCases
	authManager  util.AuthManager
	name         string
	path         string
}

func NewOrderModule(cartUseCases usecases.CartUseCases, authManager util.AuthManager) Module {
	return OrderModule{
		cartUseCases: cartUseCases,
		authManager:  authManager,
		name:         "orders",
		path:         "/orders",
	}
}

func (m OrderModule) Name() string { return m.name }
func (m OrderModule) Path() string { return m.path }

func (m OrderModule) RegisterRoutes(router *mux.Router) {
	auth := middleware.AuthMiddlewareWithAdmin(m.authManager, false)

	routes := []ModuleRoute{
		{
			Name:    "GetOrders",
			Path:    "",
			Handler: auth(m.getOrders),
			Methods: []string{http.MethodGet},
		},
	}

	for _, route := range routes {
		router.HandleFunc(m.path+route.Path, route.Handler).Methods(route.Methods...)
	}
}

func (m OrderModule) getOrders(w http.ResponseWriter, r *http.Request) {
	user, ok := util.GetUserFromContext(r.Context())
	if !ok {
		util.Write(w, util.ServerResponse{
			Status:  http.StatusUnauthorized,
			Message: "Usuário não autenticado",
		})
		return
	}

	orders, err := m.cartUseCases.GetOrders(r.Context(), int64(user.UserID))
	if err != nil {
		slog.Error("failed to get orders", "cause", err)
		util.WriteInternalError(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

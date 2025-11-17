package modules

import (
	"cachacariaapi/domain/usecases"
	"cachacariaapi/infrastructure/middleware"
	"cachacariaapi/infrastructure/util"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// CartModule handles cart endpoints
type CartModule struct {
	cartUseCases usecases.CartUseCases
	authManager  util.AuthManager
	name         string
	path         string
}

func NewCartModule(cartUseCases usecases.CartUseCases, authManager util.AuthManager) Module {
	return CartModule{
		cartUseCases: cartUseCases,
		authManager:  authManager,
		name:         "cart",
		path:         "/cart",
	}
}

func (m CartModule) Name() string { return m.name }
func (m CartModule) Path() string { return m.path }

func (m CartModule) RegisterRoutes(router *mux.Router) {
	auth := middleware.AuthMiddlewareWithAdmin(m.authManager, false)

	routes := []ModuleRoute{
		{
			Name:    "AddToCart",
			Path:    "",
			Handler: auth(m.addToCart),
			Methods: []string{http.MethodPost},
		},
		{
			Name:    "GetCart",
			Path:    "",
			Handler: auth(m.getCart),
			Methods: []string{http.MethodGet},
		},
		{
			Name:    "UpdateCartItem",
			Path:    "/{product_id}",
			Handler: auth(m.updateCartItem),
			Methods: []string{http.MethodPatch},
		},
		{
			Name:    "DeleteCartItem",
			Path:    "/{product_id}",
			Handler: auth(m.deleteCartItem),
			Methods: []string{http.MethodDelete},
		},
		{
			Name:    "Buy",
			Path:    "/buy",
			Handler: auth(m.buy),
			Methods: []string{http.MethodPost},
		},
		{
			Name:    "ClearCart",
			Path:    "/clear",
			Handler: auth(m.clearCart),
			Methods: []string{http.MethodDelete},
		},
	}

	for _, route := range routes {
		router.HandleFunc(m.path+route.Path, route.Handler).Methods(route.Methods...)
	}
}

func (m CartModule) addToCart(w http.ResponseWriter, r *http.Request) {
	user, ok := util.GetUserFromContext(r.Context())
	if !ok {
		util.Write(w, util.ServerResponse{Status: http.StatusUnauthorized, Message: "Usuário não autenticado"})
		return
	}

	var req struct {
		ProductID int64 `json:"product_id"`
		Quantity  int   `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.WriteBadRequest(w)
		return
	}

	if req.Quantity <= 0 {
		req.Quantity = 1
	}

	status, err := m.cartUseCases.AddToCart(r.Context(), int64(user.UserID), req.ProductID, req.Quantity)
	if err != nil {
		slog.Error("failed to add products to cart", "cause", err)
		util.WriteInternalError(w)
		return
	}

	util.Write(w, util.ServerResponse{
		Status:  status.Int(),
		Message: status.String(),
	})
}

func (m CartModule) getCart(w http.ResponseWriter, r *http.Request) {
	user, ok := util.GetUserFromContext(r.Context())
	if !ok {
		util.Write(w, util.ServerResponse{Status: http.StatusUnauthorized, Message: "Usuário não autenticado"})
		return
	}

	items, err := m.cartUseCases.GetCartItems(r.Context(), int64(user.UserID))
	if err != nil {
		slog.Error("failed to get cart items", "cause", err)
		util.WriteInternalError(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func (m CartModule) updateCartItem(w http.ResponseWriter, r *http.Request) {
	user, ok := util.GetUserFromContext(r.Context())
	if !ok {
		util.Write(w, util.ServerResponse{Status: http.StatusUnauthorized, Message: "Usuário não autenticado"})
		return
	}

	productIDStr := mux.Vars(r)["product_id"]
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		util.WriteBadRequest(w)
		return
	}

	var req struct {
		Quantity int `json:"quantity"`
	}
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.WriteBadRequest(w)
		return
	}

	err = m.cartUseCases.UpdateCartItem(r.Context(), int64(user.UserID), productID, req.Quantity)
	if err != nil {
		slog.Error("failed to update cart item", "cause", err)
		util.WriteInternalError(w)
		return
	}

	util.Write(w, util.ServerResponse{Status: http.StatusOK, Message: "Quantidade atualizada com sucesso"})
}

func (m CartModule) deleteCartItem(w http.ResponseWriter, r *http.Request) {
	user, ok := util.GetUserFromContext(r.Context())
	if !ok {
		util.Write(w, util.ServerResponse{Status: http.StatusUnauthorized, Message: "Usuário não autenticado"})
		return
	}

	productIDStr := mux.Vars(r)["product_id"]
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		util.WriteBadRequest(w)
		return
	}

	err = m.cartUseCases.DeleteCartItem(r.Context(), int64(user.UserID), productID)
	if err != nil {
		slog.Error("failed to delete cart item", "cause", err)
		util.WriteInternalError(w)
		return
	}

	util.Write(w, util.ServerResponse{Status: http.StatusOK, Message: "Produto removido do carrinho"})
}

func (m CartModule) buy(w http.ResponseWriter, r *http.Request) {
	user, ok := util.GetUserFromContext(r.Context())
	if !ok {
		util.Write(w, util.ServerResponse{Status: http.StatusUnauthorized, Message: "Usuário não autenticado"})
		return
	}

	status, err := m.cartUseCases.BuyItems(r.Context(), int64(user.UserID))
	if err != nil {
		slog.Error("failed to clear cart", "cause", err)
		util.WriteInternalError(w)
		return
	}

	util.Write(w, util.ServerResponse{
		Status:  status.Int(),
		Message: status.String(),
	})
}

func (m CartModule) clearCart(w http.ResponseWriter, r *http.Request) {
	user, ok := util.GetUserFromContext(r.Context())
	if !ok {
		util.Write(w, util.ServerResponse{Status: http.StatusUnauthorized, Message: "Usuário não autenticado"})
		return
	}

	err := m.cartUseCases.ClearCart(r.Context(), int64(user.UserID))
	if err != nil {
		slog.Error("failed to clear cart", "cause", err)
		util.WriteInternalError(w)
		return
	}

	util.Write(w, util.ServerResponse{Status: http.StatusOK, Message: "Carrinho limpo com sucesso"})
}

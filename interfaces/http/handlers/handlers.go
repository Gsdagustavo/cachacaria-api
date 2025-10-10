package handlers

import (
	"cachacariaapi/infrastructure/config"
	"cachacariaapi/infrastructure/modules"
	"log/slog"

	"github.com/gorilla/mux"
)

// === ROUTER ===
type MuxRouter struct {
	router *mux.Router
	cfg    *config.Server
}

func NewMuxRouter(cfg *config.Server) *MuxRouter {
	return &MuxRouter{router: mux.NewRouter(), cfg: cfg}
}

func (r *MuxRouter) StartServer(serverConfig *config.Server) {
	slog.Info("Server configuration", "port", serverConfig.Port, "host", serverConfig.Host, "baseURL", serverConfig.BaseURL)

	// Middlewares
	r.router.Use(LoggingMiddleware)
	mux.CORSMethodMiddleware(r.router)

	r.cfg.Router.PathPrefix("/api").Subrouter()
	err := r.cfg.Server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func (r *MuxRouter) SetupRoutes(modules []modules.Module) {
	for _, module := range modules {
		module.RegisterRoutes(r.router)
	}
}

//func (r *MuxRouter) registerHandlers(h *Handlers) {
//
//	router := r.router.PathPrefix("/api").Subrouter()
//
//	// This serves all files from /app/images as /images/*
//	router.PathPrefix("/images/").Handler(http.StripPrefix("/images/", http.FileServer(http.Dir("/app/images"))))
//
//	// user related handlers
//	router.HandleFunc("/users", Handle(AuthMiddleware(h.UserHandler.GetAll)))
//	router.HandleFunc("/users/{id}", Handle(AuthMiddleware(h.UserHandler.GetUser)))
//	router.HandleFunc("/users/delete/{id}", Handle(AuthMiddleware(h.UserHandler.DeleteUser)))
//	router.HandleFunc("/users/update/{id}", Handle(AuthMiddleware(h.UserHandler.UpdateUser)))
//
//	// auth handlers
//	//router.HandleFunc("/auth/register", Handle(h.AuthHandler.Register))
//	//router.HandleFunc("/auth/login", Handle(h.AuthHandler.Login))
//
//	h.AuthHandler.RegisterRoutes(router)
//
//	// product handlers
//	router.HandleFunc("/products/add", Handle(AuthMiddleware(h.ProductHandler.Add)))
//	router.HandleFunc("/products", Handle(h.ProductHandler.GetAll))
//	router.HandleFunc("/products/{id}", Handle(h.ProductHandler.GetProduct))
//	router.HandleFunc("/products/delete/{id}", Handle(AuthMiddleware(h.ProductHandler.DeleteProduct)))
//	router.HandleFunc("/products/update/{id}", Handle(AuthMiddleware(h.ProductHandler.UpdateProduct)))
//}

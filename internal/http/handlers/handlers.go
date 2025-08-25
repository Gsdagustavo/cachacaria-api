package handlers

import (
	"cachacariaapi/internal/http/core"
	"cachacariaapi/internal/http/handlers/authhandler"
	"cachacariaapi/internal/http/handlers/userhandler"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// === ROUTER ===
type MuxRouter struct {
	router *mux.Router
}

func NewMuxRouter() *MuxRouter {
	return &MuxRouter{router: mux.NewRouter()}
}

func (r *MuxRouter) StartServer(h Handlers, port string) {
	r.registerHandlers(h)
	r.serveHTTP(port)
}

func (r *MuxRouter) serveHTTP(port string) {
	log.Printf("Server is listening on port %v", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), r.router)

	if err != nil {
		log.Fatal(err)
	}
}

// === HANDLERS ===
type Handlers struct {
	UserHandler *userhandler.UserHandler
	AuthHandler *authhandler.AuthHandler
}

func (r *MuxRouter) registerHandlers(h Handlers) {
	// user related handlers
	r.router.HandleFunc("/users", Handle(h.UserHandler.GetAll))
	r.router.HandleFunc("/users/id", Handle(h.UserHandler.GetUser))
	r.router.HandleFunc("/users/add", Handle(h.UserHandler.AddUser))
	r.router.HandleFunc("/users/delete", Handle(h.UserHandler.DeleteUser))
	r.router.HandleFunc("/users/update", Handle(h.UserHandler.UpdateUser))

	// auth handlers
	r.router.HandleFunc("/auth/register", Handle(h.AuthHandler.Register))
	r.router.HandleFunc("/auth/login", Handle(h.AuthHandler.Login))

	// docs
	r.router.HandleFunc("/docs", func(w http.ResponseWriter, req *http.Request) {
		http.ServeFile(w, req, "index.html")
	})
}

type HandlerFunc func(http.ResponseWriter, *http.Request) *core.ApiError

func Handle(h HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if apiErr := h(w, r); apiErr != nil {
			apiErr.WriteHTTP(w)
		}
	}
}

package handlers

import (
	"cachacariaapi/internal/handlers/user"
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
	UserHandler *user.UserHandler
}

func (r *MuxRouter) registerHandlers(h Handlers) {

	// user related handlers
	r.router.HandleFunc("/users", h.UserHandler.GetUsers)
	r.router.HandleFunc("/users/id", h.UserHandler.GetUser)
	r.router.HandleFunc("/users/add", h.UserHandler.AddUser)
	r.router.HandleFunc("/users/delete", h.UserHandler.DeleteUser)

	// todo: products related handlers
}

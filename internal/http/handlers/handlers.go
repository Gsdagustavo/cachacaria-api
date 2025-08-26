package handlers

import (
	"cachacariaapi/internal/http/core"
	"cachacariaapi/internal/http/handlers/authhandler"
	"cachacariaapi/internal/http/handlers/userhandler"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
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
	//r.router.HandleFunc("/users/add", Handle(h.UserHandler.AddUser))
	r.router.HandleFunc("/users/delete", Handle(h.UserHandler.DeleteUser))
	r.router.HandleFunc("/users/update", Handle(h.UserHandler.UpdateUser))

	// auth handlers
	r.router.HandleFunc("/auth/register", Handle(h.AuthHandler.Register))
	r.router.HandleFunc("/auth/login", Handle(h.AuthHandler.Login))

	r.router.HandleFunc("/docs", Handle(AuthMiddleware(func(w http.ResponseWriter, req *http.Request) *core.ApiError {
		http.ServeFile(w, req, "index.html")
		log.Printf("user enterd in docs middleware")
		return nil
	})))
}

type HandlerFunc func(http.ResponseWriter, *http.Request) *core.ApiError

func Handle(h HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if apiErr := h(w, r); apiErr != nil {
			apiErr.WriteHTTP(w)
		}
	}
}

func AuthMiddleware(next HandlerFunc) HandlerFunc {
	log.Printf("on auth middleware")

	return func(w http.ResponseWriter, r *http.Request) *core.ApiError {
		cookie, err := r.Cookie("auth_token")
		if err != nil {
			return &core.ApiError{
				Code:    http.StatusUnauthorized,
				Message: "unauthorized: no token provided",
				Err:     err,
			}
		}

		tokenStr := cookie.Value
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// TODO: move the secret key to a .env file
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			return &core.ApiError{
				Code:    http.StatusUnauthorized,
				Message: "unauthorized: invalid token",
				Err:     err,
			}
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			return &core.ApiError{
				Code:    http.StatusUnauthorized,
				Message: "unauthorized: invalid token",
				Err:     err,
			}
		}

		ctx := context.WithValue(r.Context(), "user_id", claims["user_id"])
		return next(w, r.WithContext(ctx))
	}
}

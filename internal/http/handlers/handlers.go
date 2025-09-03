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
	"strings"

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
	r.router.Use(CORSMiddleware)
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
	return func(w http.ResponseWriter, r *http.Request) *core.ApiError {
		var tokenString string

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" && strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			cookie, err := r.Cookie("auth_token")
			if err != nil {
				return &core.ApiError{
					Code:    http.StatusUnauthorized,
					Message: "unauthorized: no token provided",
					Err:     err,
				}
			}

			tokenString = cookie.Value
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

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

		userID, ok := claims["user_id"].(float64)
		if !ok {
			return &core.ApiError{
				Code:    http.StatusUnauthorized,
				Message: "unauthorized: invalid user_id",
			}
		}

		ctx := context.WithValue(r.Context(), "user_id", int64(userID))
		return next(w, r.WithContext(ctx))
	}
}

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		origins := map[string]bool{
			"http://127.0.0.1:5000": true,
			"http://127.0.0.1:5001": true,
			"http://127.0.0.1:5500": true,
			"http://127.0.0.1:5501": true,
		}

		if origins[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Expose-Headers", "Authorization")

		if r.Method == http.MethodOptions {
			log.Printf("[CORS Middleware] allow options | no content")
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

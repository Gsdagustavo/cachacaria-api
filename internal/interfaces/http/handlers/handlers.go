package handlers

import (
	"cachacariaapi/internal/infrastructure/config"
	"cachacariaapi/internal/interfaces/http/core"
	"cachacariaapi/internal/interfaces/http/handlers/authhandler"
	"cachacariaapi/internal/interfaces/http/handlers/producthandler"
	"cachacariaapi/internal/interfaces/http/handlers/userhandler"
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

// === ROUTER ===
type MuxRouter struct {
	router *mux.Router
	cfg    *config.ServerConfig
}

func NewMuxRouter(cfg *config.ServerConfig) *MuxRouter {
	return &MuxRouter{router: mux.NewRouter(), cfg: cfg}
}

func (r *MuxRouter) StartServer(h *Handlers, serverConfig *config.ServerConfig) {
	slog.Info("Server configuration", "port", serverConfig.Port, "address", serverConfig.Address, "baseURL", serverConfig.BaseURL)

	r.registerHandlers(h)
	r.router.Use(CORSMiddleware)
	r.serveHTTP(serverConfig.Address + ":" + serverConfig.Port)
}

func (r *MuxRouter) serveHTTP(serverAddress string) {
	err := http.ListenAndServe(serverAddress, r.router)
	log.Printf("Server is listening on %s", serverAddress)
	if err != nil {
		log.Fatal(err)
	}
}

// === HANDLERS ===
type Handlers struct {
	UserHandler    *userhandler.UserHandler
	AuthHandler    *authhandler.AuthHandler
	ProductHandler *producthandler.ProductHandler
}

func NewHandlers(userHandler *userhandler.UserHandler, authHandler *authhandler.AuthHandler, productHandler *producthandler.ProductHandler) *Handlers {
	return &Handlers{UserHandler: userHandler, AuthHandler: authHandler, ProductHandler: productHandler}
}

func (r *MuxRouter) registerHandlers(h *Handlers) {
	// This serves all files from /app/images as /images/*
	r.router.PathPrefix("/images/").Handler(http.StripPrefix("/images/",
		http.FileServer(http.Dir("/app/images")),
	))

	r.router.Use(CORSMiddleware)
	r.router.Use(LoggingMiddleware)

	// user related handlers
	r.router.HandleFunc("/users", Handle(h.UserHandler.GetAll))
	r.router.HandleFunc("/users/{id}", Handle(h.UserHandler.GetUser))
	r.router.HandleFunc("/users/delete/{id}", Handle(h.UserHandler.DeleteUser))
	r.router.HandleFunc("/users/update/{id}", Handle(h.UserHandler.UpdateUser))

	// auth handlers
	r.router.HandleFunc("/auth/register", Handle(h.AuthHandler.Register))
	r.router.HandleFunc("/auth/login", Handle(h.AuthHandler.Login))

	// product handlers
	r.router.HandleFunc("/products/add", Handle(h.ProductHandler.Add))
	r.router.HandleFunc("/products", Handle(h.ProductHandler.GetAll))
	r.router.HandleFunc("/products/{id}", Handle(h.ProductHandler.GetProduct))
	r.router.HandleFunc("/products/delete/{id}", Handle(h.ProductHandler.DeleteProduct))
	r.router.HandleFunc("/products/update/{id}", Handle(h.ProductHandler.UpdateProduct))

	r.router.HandleFunc("/docs", Handle(AuthMiddleware(func(w http.ResponseWriter, req *http.Request) *core.ServerError {
		http.ServeFile(w, req, "index.html")
		log.Printf("user enterd in docs middleware")
		return nil
	})))
}

type HandlerFunc func(http.ResponseWriter, *http.Request) *core.ServerError

func Handle(h HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if apiErr := h(w, r); apiErr != nil {
			apiErr.WriteHTTP(w)
		}
	}
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Incoming request",
			"method", r.Method,
			"url", r.URL.Path,
			"remote_addr", r.RemoteAddr,
			"user_agent", r.UserAgent(),
			"host", r.Host,
			"cookies", r.Cookies(),
			"body", r.Body,
			"form", r.Form,
			"post_form", r.PostForm,
			"multipart_form", r.MultipartForm,
		)
		next.ServeHTTP(w, r)
	})
}

func AuthMiddleware(next HandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *core.ServerError {
		var tokenString string

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" && strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			cookie, err := r.Cookie("auth_token")
			if err != nil {
				return &core.ServerError{
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
			return &core.ServerError{
				Code:    http.StatusUnauthorized,
				Message: "unauthorized: invalid token",
				Err:     err,
			}
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			return &core.ServerError{
				Code:    http.StatusUnauthorized,
				Message: "unauthorized: invalid token",
				Err:     err,
			}
		}

		userID, ok := claims["user_id"].(float64)
		if !ok {
			return &core.ServerError{
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
			"http://localhost:5173": true,
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

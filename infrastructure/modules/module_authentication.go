package modules

import (
	"cachacariaapi/domain/entities"
	"cachacariaapi/domain/usecases"
	"cachacariaapi/infrastructure/util"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type AuthResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
}

type AuthModule struct {
	authUseCases *usecases.AuthUseCases
	name         string
	path         string
}

func NewAuthModule(authUseCases *usecases.AuthUseCases) *AuthModule {
	return &AuthModule{
		authUseCases: authUseCases,
		name:         "auth",
		path:         "/auth",
	}
}

func (a AuthModule) Name() string {
	return a.name
}

func (a AuthModule) Path() string {
	return a.path
}

func (a AuthModule) RegisterRoutes(router *mux.Router) {
	routes := []ModuleRoute{
		{
			Name:    "Login",
			Path:    "/login",
			Handler: a.login,
			Methods: []string{http.MethodPost},
		},
		{
			Name:    "Register",
			Path:    "/register",
			Handler: a.register,
			Methods: []string{http.MethodPost},
		},
		{
			Name:    "Return user info",
			Path:    "/me",
			Handler: a.getData,
			Methods: []string{http.MethodPost},
		},
	}

	for _, route := range routes {
		router.HandleFunc(a.path+route.Path, route.Handler).Methods(route.Methods...)
	}
}

func (a AuthModule) login(w http.ResponseWriter, r *http.Request) {
	var credentials entities.UserCredentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		util.WriteBadRequest(w)
		return
	}

	token, statusCode, err := a.authUseCases.AttemptLogin(r.Context(), credentials)
	if err != nil {
		slog.Error("failed to login", "cause", err)
		util.WriteInternalError(w)
		return
	}

	response := AuthResponse{
		Status:  statusCode.Int(),
		Message: statusCode.String(),
		Token:   token,
	}

	util.Write(w, response)
}

func (a AuthModule) register(w http.ResponseWriter, r *http.Request) {
	var credentials entities.UserCredentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		util.WriteBadRequest(w)
		return
	}

	statusCode, err := a.authUseCases.RegisterUser(r.Context(), credentials)
	if err != nil {
		slog.Error("failed to register user", "cause", err)
		util.WriteInternalError(w)
		return
	}

	response := util.ServerResponse{
		Status:  statusCode.Int(),
		Message: statusCode.String(),
	}

	util.Write(w, response)
}

func (a AuthModule) getData(w http.ResponseWriter, r *http.Request) {
	token := util.GetAuthTokenFromRequest(r)

	user, err := a.authUseCases.GetUserByAuthToken(token)
	if err != nil {
		util.WriteResponse(w, util.ServerResponse{
			Status:  http.StatusUnauthorized,
			Message: "Token inválido ou expirado",
		})
		return
	}

	if user == nil {
		util.WriteResponse(w, util.ServerResponse{
			Status:  http.StatusUnauthorized,
			Message: "Token inválido ou expirado",
		})
		return
	}

	util.Write(w, user)
}

func (a AuthModule) SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader != "" {
			token := strings.TrimPrefix(authHeader, "Bearer ")
			user, err := a.authUseCases.GetUserByAuthToken(token)

			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			ctx := util.NewContextWithUser(r.Context(), user.ID, user.IsAdm)

			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

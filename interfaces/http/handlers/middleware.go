package handlers

import (
	"cachacariaapi/interfaces/http/core"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Incoming request", "method", r.Method, "url", r.URL.Path, "remote_addr", r.RemoteAddr, "user_agent", r.UserAgent(), "host", r.Host, "cookies", r.Cookies(), "body", r.Body, "form", r.Form, "post_form", r.PostForm, "multipart_form", r.MultipartForm)
		next.ServeHTTP(w, r)
	})
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var tokenString string

		authHeader := r.Header.Get("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			cookie, err := r.Cookie("auth_token")
			if err != nil {
				core.WriteServerError(w, &core.ServerError{
					Code:    http.StatusUnauthorized,
					Message: "no token provided",
					Err:     err,
				})
				return
			}

			tokenString = cookie.Value
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			jwtSecret := os.Getenv("JWT_SECRET")
			if jwtSecret == "" {
				return nil, fmt.Errorf("JWT_SECRET environment variable not set")
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			core.WriteServerError(w, &core.ServerError{
				Code:    http.StatusUnauthorized,
				Message: "invalid token",
				Err:     err,
			})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			core.WriteServerError(w, &core.ServerError{
				Code:    http.StatusUnauthorized,
				Message: "invalid token structure",
				Err:     err,
			})
			return
		}

		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			core.WriteServerError(w, &core.ServerError{
				Code:    http.StatusUnauthorized,
				Message: "invalid user_id format in token",
			})
			return
		}

		userID := int(userIDFloat)

		ctx := r.Context()
		ctx = core.NewContextWithUserID(ctx, userID)

		r = r.WithContext(ctx)

		next(w, r)
	}
}

package middleware

import (
	"cachacariaapi/infrastructure/util"
	"log/slog"
	"net/http"
	"strings"
)

func AuthMiddlewareWithAdmin(authManager util.AuthManager, adminOnly bool) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				slog.Info("No Authorization header")
				util.WriteResponse(w, util.ServerResponse{
					Status:  http.StatusUnauthorized,
					Message: "Sem autorização",
				})
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				slog.Info("Invalid Authorization header format")
				util.WriteResponse(w, util.ServerResponse{
					Status:  http.StatusUnauthorized,
					Message: "Sem autorização",
				})
				return
			}

			token := parts[1]
			payload, err := authManager.VerifyToken(token)
			if err != nil {
				slog.Info("Invalid token")
				util.WriteResponse(w, util.ServerResponse{
					Status:  http.StatusUnauthorized,
					Message: "Sem autorização",
				})
				return
			}

			if adminOnly && !payload.IsAdmin {
				slog.Info("Not admin")
				util.WriteResponse(w, util.ServerResponse{
					Status:  http.StatusUnauthorized,
					Message: "Sem autorização",
				})
				return
			}

			ctx := util.NewContextWithUser(r.Context(), payload.UserID, payload.IsAdmin)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}

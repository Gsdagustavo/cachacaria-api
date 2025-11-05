package middleware

import (
	"cachacariaapi/infrastructure/util"
	"net/http"
	"strings"
)

// AuthMiddleware validates the PASETO token in the Authorization header
func AuthMiddleware(crypt util.Crypt) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				util.Write(w, util.ServerResponse{
					Status:  http.StatusUnauthorized,
					Message: "Token ausente no cabeçalho Authorization",
				})
				return
			}

			// Expected format: "Bearer <token>"
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				util.Write(w, util.ServerResponse{
					Status:  http.StatusUnauthorized,
					Message: "Formato de token inválido. Use 'Bearer <token>'",
				})
				return
			}

			token := parts[1]
			payload, err := crypt.VerifyAuthToken(token)
			if err != nil {
				util.Write(w, util.ServerResponse{
					Status:  http.StatusUnauthorized,
					Message: "Token inválido ou expirado",
				})
				return
			}

			// Store userID in the request context for use in handlers
			ctx := util.NewContextWithUserID(r.Context(), payload.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}

package middleware

import (
	"cachacariaapi/infrastructure/util"
	"net/http"
	"strings"
)

func AuthMiddlewareWithAdmin(crypt util.Crypt, adminOnly bool) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				util.WriteResponse(w, util.ServerResponse{
					Status:  http.StatusUnauthorized,
					Message: "Sem autorização",
				})
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				util.WriteResponse(w, util.ServerResponse{
					Status:  http.StatusUnauthorized,
					Message: "Sem autorização",
				})
				return
			}

			token := parts[1]
			payload, err := crypt.VerifyAuthToken(token)
			if err != nil {
				util.WriteResponse(w, util.ServerResponse{
					Status:  http.StatusUnauthorized,
					Message: "Sem autorização",
				})
				return
			}

			if adminOnly && !payload.IsAdmin {
				util.WriteResponse(w, util.ServerResponse{
					Status:  http.StatusUnauthorized,
					Message: "Sem autorização",
				})
				return
			}

			// Armazena o ID do usuário no contexto
			ctx := util.NewContextWithUser(r.Context(), payload.UserID, payload.IsAdmin)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}

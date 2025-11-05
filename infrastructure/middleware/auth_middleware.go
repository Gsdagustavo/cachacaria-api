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
				util.Write(w, util.ServerResponse{
					Status:  http.StatusUnauthorized,
					Message: "Token ausente no cabeçalho Authorization",
				})
				return
			}

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

			if adminOnly && !payload.IsAdmin {
				util.Write(w, util.ServerResponse{
					Status:  http.StatusForbidden,
					Message: "Acesso negado: apenas administradores podem acessar esta rota",
				})
				return
			}

			// Armazena o ID do usuário no contexto
			ctx := util.NewContextWithUser(r.Context(), payload.UserID, payload.IsAdmin)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}

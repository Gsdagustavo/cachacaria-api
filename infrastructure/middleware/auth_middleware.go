package middleware

import (
	"cachacariaapi/infrastructure/util"
	"log"
	"net/http"
	"strings"
)

func AuthMiddlewareWithAdmin(crypt util.Crypt, adminOnly bool) func(http.HandlerFunc) http.HandlerFunc {
	log.Printf("AuthMiddlewareWithAdmin")

	return func(next http.HandlerFunc) http.HandlerFunc {

		log.Printf("AuthMiddlewareWithAdmin 1st")

		return func(w http.ResponseWriter, r *http.Request) {
			log.Printf("AuthMiddlewareWithAdmin 2nd")

			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				log.Printf("no token found in request")
				util.Write(w, util.ServerResponse{
					Status:  http.StatusUnauthorized,
					Message: "Token ausente no cabeçalho Authorization",
				})
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				log.Printf("invalid token found in request")
				util.Write(w, util.ServerResponse{
					Status:  http.StatusUnauthorized,
					Message: "Formato de token inválido. Use 'Bearer <token>'",
				})
				return
			}

			token := parts[1]
			payload, err := crypt.VerifyAuthToken(token)
			if err != nil {
				log.Printf("invalid/expired token found in request")
				util.Write(w, util.ServerResponse{
					Status:  http.StatusUnauthorized,
					Message: "Token inválido ou expirado",
				})
				return
			}

			log.Printf("admin only: %v", adminOnly)
			log.Printf("is admin: %v", payload.IsAdmin)

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

package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/akinolaemmanuel49/gocommerce/internal/auth/models"
	"github.com/akinolaemmanuel49/gocommerce/internal/auth/utils"
)

type contextKey string

const UserClaimsKey contextKey = "userClaims"

// AuthMiddleware validates JWT tokens and adds claims to the request context
func AuthMiddleware(jwtSecretKey []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract the token from the Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
				return
			}

			// Token format: "Bearer <token>"
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
				return
			}

			// Parse and validate the token
			claims, err := utils.ParseJWT(jwtSecretKey, tokenString)
			if err != nil {
				http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
				return
			}

			// Add the claims to the request context
			ctx := context.WithValue(r.Context(), UserClaimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetClaims retrieves JWT claims from the request context
func GetClaims(ctx context.Context) *models.JWTClaims {
	if claims, ok := ctx.Value(UserClaimsKey).(*models.JWTClaims); ok {
		return claims
	}
	return nil
}

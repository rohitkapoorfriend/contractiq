package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/contractiq/contractiq/internal/infrastructure/auth"
	"github.com/contractiq/contractiq/internal/interfaces/http/response"
	"github.com/contractiq/contractiq/pkg/apperror"
)

const UserIDKey contextKey = "user_id"
const UserEmailKey contextKey = "user_email"

// Auth verifies the JWT token in the Authorization header.
func Auth(jwtService *auth.JWTService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" {
				response.Error(w, apperror.NewUnauthorized("missing authorization header"))
				return
			}

			parts := strings.SplitN(header, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
				response.Error(w, apperror.NewUnauthorized("invalid authorization format"))
				return
			}

			claims, err := jwtService.ValidateToken(parts[1])
			if err != nil {
				response.Error(w, err)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, UserEmailKey, claims.Email)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserID extracts the authenticated user ID from the context.
func GetUserID(ctx context.Context) string {
	id, _ := ctx.Value(UserIDKey).(string)
	return id
}

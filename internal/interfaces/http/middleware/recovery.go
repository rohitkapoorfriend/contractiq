package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/contractiq/contractiq/internal/interfaces/http/response"
	"github.com/contractiq/contractiq/pkg/apperror"
	"go.uber.org/zap"
)

// Recovery catches panics and returns a 500 error instead of crashing.
func Recovery(logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					logger.Error("panic recovered",
						zap.Any("panic", rec),
						zap.String("stack", string(debug.Stack())),
						zap.String("request_id", GetRequestID(r.Context())),
					)
					response.Error(w, apperror.NewInternal(nil))
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

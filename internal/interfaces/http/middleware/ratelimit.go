package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/contractiq/contractiq/internal/interfaces/http/response"
	"github.com/contractiq/contractiq/pkg/apperror"
)

type visitor struct {
	tokens    float64
	lastSeen  time.Time
}

// RateLimit implements a simple token bucket rate limiter per IP.
func RateLimit(rps float64, burst int) func(next http.Handler) http.Handler {
	var (
		mu       sync.Mutex
		visitors = make(map[string]*visitor)
	)

	// Clean up old visitors periodically.
	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock()
			for ip, v := range visitors {
				if time.Since(v.lastSeen) > 3*time.Minute {
					delete(visitors, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr

			mu.Lock()
			v, exists := visitors[ip]
			if !exists {
				v = &visitor{tokens: float64(burst)}
				visitors[ip] = v
			}

			// Refill tokens based on time elapsed.
			elapsed := time.Since(v.lastSeen).Seconds()
			v.tokens += elapsed * rps
			if v.tokens > float64(burst) {
				v.tokens = float64(burst)
			}
			v.lastSeen = time.Now()

			if v.tokens < 1 {
				mu.Unlock()
				w.Header().Set("Retry-After", "1")
				response.JSON(w, http.StatusTooManyRequests, response.ErrorResponse{
					Error: response.ErrorBody{
						Code:    string(apperror.CodeBadRequest),
						Message: "rate limit exceeded",
					},
				})
				return
			}

			v.tokens--
			mu.Unlock()

			next.ServeHTTP(w, r)
		})
	}
}

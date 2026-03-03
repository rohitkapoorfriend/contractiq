package http

import (
	"strings"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/contractiq/contractiq/internal/infrastructure/auth"
	"github.com/contractiq/contractiq/internal/interfaces/http/handler"
	"github.com/contractiq/contractiq/internal/interfaces/http/middleware"
	"go.uber.org/zap"
)

// RouterDeps contains all dependencies needed to build the router.
type RouterDeps struct {
	Logger          *zap.Logger
	JWTService      *auth.JWTService
	ContractHandler *handler.ContractHandler
	TemplateHandler *handler.TemplateHandler
	PartyHandler    *handler.PartyHandler
	AuthHandler     *handler.AuthHandler
	HealthHandler   *handler.HealthHandler
	AllowedOrigins  string
}

// NewRouter creates the Chi router with all middleware and routes.
func NewRouter(deps RouterDeps) *chi.Mux {
	r := chi.NewRouter()

	// Global middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.Recovery(deps.Logger))
	r.Use(middleware.Logging(deps.Logger))
	r.Use(chimiddleware.RealIP)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   strings.Split(deps.AllowedOrigins, ","),
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Request-ID"},
		ExposedHeaders:   []string{"X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	r.Use(middleware.RateLimit(100, 200))

	r.Route("/api/v1", func(r chi.Router) {
		// Public endpoints
		r.Get("/health", deps.HealthHandler.Check)

		// Auth endpoints (public)
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", deps.AuthHandler.Register)
			r.Post("/login", deps.AuthHandler.Login)
		})

		// Protected endpoints
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(deps.JWTService))

			// Contracts
			r.Route("/contracts", func(r chi.Router) {
				r.Post("/", deps.ContractHandler.Create)
				r.Get("/", deps.ContractHandler.List)
				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", deps.ContractHandler.Get)
					r.Put("/", deps.ContractHandler.Update)
					r.Post("/submit", deps.ContractHandler.Submit)
					r.Post("/approve", deps.ContractHandler.Approve)
					r.Post("/sign", deps.ContractHandler.Sign)
					r.Post("/terminate", deps.ContractHandler.Terminate)
				})
			})

			// Templates
			r.Route("/templates", func(r chi.Router) {
				r.Post("/", deps.TemplateHandler.Create)
				r.Get("/", deps.TemplateHandler.List)
				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", deps.TemplateHandler.Get)
					r.Put("/", deps.TemplateHandler.Update)
					r.Delete("/", deps.TemplateHandler.Delete)
				})
			})

			// Parties
			r.Route("/parties", func(r chi.Router) {
				r.Post("/", deps.PartyHandler.Create)
				r.Get("/", deps.PartyHandler.List)
				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", deps.PartyHandler.Get)
					r.Put("/", deps.PartyHandler.Update)
					r.Delete("/", deps.PartyHandler.Delete)
				})
			})
		})
	})

	return r
}

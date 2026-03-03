package handler

import (
	"net/http"

	"github.com/contractiq/contractiq/internal/interfaces/http/response"
	"github.com/jmoiron/sqlx"
)

// HealthHandler handles the health check endpoint.
type HealthHandler struct {
	db *sqlx.DB
}

// NewHealthHandler creates a new health handler.
func NewHealthHandler(db *sqlx.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

// Check handles GET /api/v1/health.
func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	status := "ok"
	dbStatus := "up"

	if err := h.db.PingContext(r.Context()); err != nil {
		dbStatus = "down"
		status = "degraded"
	}

	response.OK(w, map[string]string{
		"status":   status,
		"database": dbStatus,
	})
}

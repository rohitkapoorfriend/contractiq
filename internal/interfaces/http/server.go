package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// Server wraps the HTTP server with graceful shutdown support.
type Server struct {
	httpServer *http.Server
	logger     *zap.Logger
}

// NewServer creates a new HTTP server.
func NewServer(handler http.Handler, host, port string, readTimeout, writeTimeout time.Duration, logger *zap.Logger) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         fmt.Sprintf("%s:%s", host, port),
			Handler:      handler,
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
		},
		logger: logger,
	}
}

// Start begins listening for HTTP requests.
func (s *Server) Start() error {
	s.logger.Info("starting HTTP server", zap.String("addr", s.httpServer.Addr))
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("HTTP server error: %w", err)
	}
	return nil
}

// Shutdown gracefully stops the server.
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("shutting down HTTP server")
	return s.httpServer.Shutdown(ctx)
}

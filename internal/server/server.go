package server

import (
	"net/http"
	"os"

	"github.com/sunwolfengineering/nyla-core/internal/middleware"
	"github.com/sunwolfengineering/nyla-core/internal/storage"
	"github.com/sunwolfengineering/nyla-core/pkg/handlers"
)

// Server represents the unified HTTP server
type Server struct {
	db *storage.DB
	mux    *http.ServeMux
	handler http.Handler
}

// New creates a new unified server instance
func New(db *storage.DB) *Server {
	s := &Server{
		db: db,
		mux:    http.NewServeMux(),
	}
	
	s.setupRoutes()
	s.setupMiddleware()
	
	return s
}

// setupRoutes configures all API and UI routes
func (s *Server) setupRoutes() {
	// Initialize handlers
	apiHandlers := &handlers.Handlers{DB: s.db}
	
	apiBaseURL := os.Getenv("API_BASE_URL")
	if apiBaseURL == "" {
		apiBaseURL = "https://api.localhost"
	}
	uiHandlers := &handlers.UIHandlers{APIBaseURL: apiBaseURL}
	
	// API routes at /api/v1/*
	s.mux.HandleFunc("GET /api/v1/collect", apiHandlers.GetCollectV1)
	s.mux.HandleFunc("GET /api/v1/stats/realtime", apiHandlers.GetStatsRealtimeV1)
	
	// UI routes
	s.mux.HandleFunc("GET /", uiHandlers.DashboardHandler)
}

// setupMiddleware configures middleware stack
func (s *Server) setupMiddleware() {
	corsConfig := middleware.NewCORSConfig()
	s.handler = corsConfig.CORS(s.mux)
}

// Handler returns the configured HTTP handler
func (s *Server) Handler() http.Handler {
	return s.handler
}

// ServeHTTP implements http.Handler interface
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.handler.ServeHTTP(w, r)
}

// ListenAndServe starts the server on the specified address
func (s *Server) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, s.handler)
}

// Package server provides the HTTP server setup and routing
package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/Krea-University/speed-test-server/docs"
	"github.com/Krea-University/speed-test-server/internal/auth"
	"github.com/Krea-University/speed-test-server/internal/config"
	"github.com/Krea-University/speed-test-server/internal/database"
	"github.com/Krea-University/speed-test-server/internal/handlers"
	"github.com/Krea-University/speed-test-server/internal/middleware"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// Server represents the HTTP server instance
type Server struct {
	httpServer *http.Server
	handlers   *handlers.Handlers
	db         *database.Service
}

// New creates a new server instance with all routes configured
func New() *Server {
	// Initialize database
	db, err := database.New()
	if err != nil {
		log.Printf("Warning: Database connection failed: %v", err)
		log.Println("Continuing without database features...")
	}

	// Initialize handlers with database
	h := handlers.New(db)

	// Initialize auth service
	var authService *auth.Service
	if db != nil {
		authService = auth.New(db)
	}

	// Create router with middleware
	r := mux.NewRouter()

	// Initialize concurrent request limiter with configurable max requests
	concurrentLimiter := middleware.NewConcurrentRequestLimiter(config.GetMaxConcurrentRequests())

	// Apply global middleware (but skip for WebSocket)
	r.Use(middleware.Logging)
	r.Use(middleware.Security)
	r.Use(middleware.CORS)

	// Create a subrouter for non-WebSocket endpoints with concurrent limiting
	api := r.PathPrefix("/").Subrouter()
	api.Use(concurrentLimiter.Middleware)

	// Apply API key authentication if database is available (to non-WebSocket endpoints)
	// Note: Rate limiting is disabled for speed test performance
	if authService != nil {
		api.Use(authService.APIKeyAuth)
	}

	// Test session endpoints (public, no auth required)
	api.HandleFunc("/test/start", h.StartTest).Methods("POST", "OPTIONS")
	api.HandleFunc("/test/complete", h.CompleteTest).Methods("POST", "OPTIONS")

	// Public speed test endpoints (with concurrent limiting)
	api.HandleFunc("/ping", h.Ping).Methods("GET", "OPTIONS")
	api.HandleFunc("/download", h.Download).Methods("GET", "OPTIONS")
	api.HandleFunc("/upload", h.Upload).Methods("POST", "OPTIONS")

	// WebSocket endpoint (without concurrent limiting to avoid hijacker issues)
	r.HandleFunc("/ws", h.WebSocket).Methods("GET", "OPTIONS")

	// Public information endpoints (with concurrent limiting)
	api.HandleFunc("/ip", h.IP).Methods("GET", "OPTIONS")
	api.HandleFunc("/healthz", h.Health).Methods("GET", "OPTIONS")
	api.HandleFunc("/version", h.Version).Methods("GET", "OPTIONS")
	api.HandleFunc("/config", h.Config).Methods("GET", "OPTIONS")

	// Web interface endpoints (public HTML pages, with concurrent limiting)
	api.HandleFunc("/", h.ServeSpeedTestHTML).Methods("GET", "HEAD", "OPTIONS")
	api.HandleFunc("/speedtest.html", h.ServeSpeedTestHTML).Methods("GET", "HEAD", "OPTIONS")
	api.HandleFunc("/new", h.ServeSpeedTestNewHTML).Methods("GET", "HEAD", "OPTIONS")

	// Ookla-compatible endpoints (public, with concurrent limiting)
	api.HandleFunc("/result/{id}", h.GetSpeedTestOokla).Methods("GET", "OPTIONS")

	// Admin endpoints (always available for monitoring, with concurrent limiting)
	admin := api.PathPrefix("/admin").Subrouter()
	admin.HandleFunc("/", h.AdminDashboard).Methods("GET")
	admin.HandleFunc("/api/stats", h.AdminStats).Methods("GET")
	admin.HandleFunc("/api/recent-tests", h.AdminRecentTests).Methods("GET")
	admin.HandleFunc("/api/system", h.AdminSystemInfo).Methods("GET")

	// API endpoints (require authentication if database is available)
	if db != nil {
		apiAuth := api.PathPrefix("/api").Subrouter()
		apiAuth.HandleFunc("/tests", h.GetAllSpeedTests).Methods("GET")
		apiAuth.HandleFunc("/tests", h.CreateSpeedTest).Methods("POST")
		apiAuth.HandleFunc("/tests/{id}", h.GetSpeedTest).Methods("GET")
	}

	// Swagger documentation endpoint
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// 404 handler - redirect to main page
	r.NotFoundHandler = http.HandlerFunc(h.NotFound)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = config.DefaultPort
	}

	// Create HTTP server
	httpServer := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	return &Server{
		httpServer: httpServer,
		handlers:   h,
		db:         db,
	}
}

// Start starts the HTTP server and handles graceful shutdown
func (s *Server) Start() error {
	// Start server in a goroutine
	go func() {
		log.Printf("Speed Test Server starting on port %s", s.httpServer.Addr)
		log.Printf("Version: %s", config.Version)
		maxConcurrent := config.GetMaxConcurrentRequests()
		if maxConcurrent == 0 {
			log.Printf("Max concurrent requests: unlimited")
		} else {
			log.Printf("Max concurrent requests: %d", maxConcurrent)
		}
		log.Printf("Available endpoints:")
		log.Printf("  GET  /ping      - Latency measurement")
		log.Printf("  GET  /download  - Download speed test")
		log.Printf("  POST /upload    - Upload speed test")
		log.Printf("  GET  /ws        - WebSocket for jitter measurement")
		log.Printf("  GET  /ip        - IP geolocation information")
		log.Printf("  GET  /healthz   - Health check")
		log.Printf("  GET  /version   - Application version")
		log.Printf("  GET  /config    - Server configuration")
		log.Printf("  GET  /speedtest.html - Main speed test interface")
		log.Printf("  GET  /new       - Modern speed test interface")
		log.Printf("  GET  /result/{id} - Ookla-compatible speed test results")
		log.Printf("  GET  /swagger/index.html - Swagger API documentation")

		if s.db != nil {
			log.Printf("API endpoints (require authentication):")
			log.Printf("  GET  /api/tests    - List all speed tests")
			log.Printf("  POST /api/tests    - Create speed test")
			log.Printf("  GET  /api/tests/{id} - Get specific speed test")
		}

		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	// Close database connection
	if s.db != nil {
		if err := s.db.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}

	log.Println("Server exited")
	return nil
}

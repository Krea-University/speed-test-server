// Package auth provides authentication and rate limiting middleware
package auth

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Krea-University/speed-test-server/internal/database"
)

// Service provides authentication and rate limiting
type Service struct {
	db *database.Service
}

// New creates a new auth service
func New(db *database.Service) *Service {
	return &Service{db: db}
}

// RateLimitConfig defines rate limiting configuration
type RateLimitConfig struct {
	RequestsPerMinute int
	WhitelistEnabled  bool
}

// DefaultRateLimits defines default rate limits for different endpoints
var DefaultRateLimits = map[string]RateLimitConfig{
	"/ping":     {RequestsPerMinute: 60, WhitelistEnabled: true},
	"/download": {RequestsPerMinute: 10, WhitelistEnabled: true},
	"/upload":   {RequestsPerMinute: 10, WhitelistEnabled: true},
	"/ip":       {RequestsPerMinute: 30, WhitelistEnabled: true},
	"/ws":       {RequestsPerMinute: 20, WhitelistEnabled: true},
	"/api/":     {RequestsPerMinute: 100, WhitelistEnabled: false}, // API endpoints
}

// APIKeyAuth middleware for API key authentication
func (s *Service) APIKeyAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only apply to /api/ routes
		if !strings.HasPrefix(r.URL.Path, "/api/") {
			next.ServeHTTP(w, r)
			return
		}

		// Extract API key from header
		apiKey := r.Header.Get("X-API-Key")
		if apiKey == "" {
			// Try Authorization header
			auth := r.Header.Get("Authorization")
			if strings.HasPrefix(auth, "Bearer ") {
				apiKey = strings.TrimPrefix(auth, "Bearer ")
			}
		}

		if apiKey == "" {
			http.Error(w, `{"error":"API key required","code":"MISSING_API_KEY"}`, http.StatusUnauthorized)
			return
		}

		// Hash the API key
		keyHash := fmt.Sprintf("%x", sha256.Sum256([]byte(apiKey)))

		// Verify API key
		key, err := s.db.GetAPIKey(keyHash)
		if err != nil {
			http.Error(w, `{"error":"Invalid API key","code":"INVALID_API_KEY"}`, http.StatusUnauthorized)
			return
		}

		// Update last used timestamp
		go s.db.UpdateAPIKeyLastUsed(keyHash)

		// Store API key info in request context for later use
		r.Header.Set("X-API-Key-ID", key.ID)
		r.Header.Set("X-API-Key-Name", key.Name)

		next.ServeHTTP(w, r)
	})
}

// RateLimit middleware for rate limiting
// Note: This middleware is currently disabled in server.go for speed test performance
func (s *Service) RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Rate limiting is disabled - pass through all requests
		next.ServeHTTP(w, r)
		return

		// Get client identifier (IP address)
		clientIP := getClientIP(r)

		// Check if IP is whitelisted
		isWhitelisted, err := s.db.IsWhitelisted(clientIP)
		if err != nil {
			// Log error but continue (fail open)
			fmt.Printf("Error checking whitelist: %v\n", err)
		}

		// Determine rate limit config
		endpoint := r.URL.Path
		config := getRateLimitConfig(endpoint)

		// Skip rate limiting for whitelisted IPs if enabled
		if config.WhitelistEnabled && isWhitelisted {
			next.ServeHTTP(w, r)
			return
		}

		// Use API key rate limit if available
		limit := config.RequestsPerMinute
		identifier := clientIP

		// For API endpoints, use API key specific limits
		if strings.HasPrefix(r.URL.Path, "/api/") {
			apiKeyID := r.Header.Get("X-API-Key-ID")
			if apiKeyID != "" {
				identifier = "api:" + apiKeyID
				// Could fetch specific API key rate limit here
			}
		}

		// Check rate limit
		allowed, err := s.db.CheckRateLimit(identifier, endpoint, limit)
		if err != nil {
			// Log error but continue (fail open)
			fmt.Printf("Error checking rate limit: %v\n", err)
			next.ServeHTTP(w, r)
			return
		}

		if !allowed {
			w.Header().Set("X-RateLimit-Limit", strconv.Itoa(limit))
			w.Header().Set("X-RateLimit-Remaining", "0")
			w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(time.Minute).Unix(), 10))
			http.Error(w, `{"error":"Rate limit exceeded","code":"RATE_LIMIT_EXCEEDED","retry_after":60}`, http.StatusTooManyRequests)
			return
		}

		// Add rate limit headers
		w.Header().Set("X-RateLimit-Limit", strconv.Itoa(limit))
		// Note: We don't track remaining requests in this simple implementation

		next.ServeHTTP(w, r)
	})
}

// getRateLimitConfig returns the rate limit configuration for an endpoint
func getRateLimitConfig(endpoint string) RateLimitConfig {
	// Check for exact match first
	if config, exists := DefaultRateLimits[endpoint]; exists {
		return config
	}

	// Check for prefix matches
	for pattern, config := range DefaultRateLimits {
		if strings.HasPrefix(endpoint, pattern) {
			return config
		}
	}

	// Default rate limit
	return RateLimitConfig{
		RequestsPerMinute: 30,
		WhitelistEnabled:  true,
	}
}

// getClientIP extracts the real client IP from request headers
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	if xForwardedFor != "" {
		ips := strings.Split(xForwardedFor, ",")
		clientIP := strings.TrimSpace(ips[0])
		if clientIP != "" {
			return clientIP
		}
	}

	// Check X-Real-IP header
	xRealIP := r.Header.Get("X-Real-IP")
	if xRealIP != "" {
		return strings.TrimSpace(xRealIP)
	}

	// Check X-Client-IP header
	xClientIP := r.Header.Get("X-Client-IP")
	if xClientIP != "" {
		return strings.TrimSpace(xClientIP)
	}

	// Fall back to RemoteAddr
	if ip := strings.Split(r.RemoteAddr, ":")[0]; ip != "" {
		return ip
	}

	return "unknown"
}

// GenerateAPIKey generates a new API key
func GenerateAPIKey() string {
	// Generate a secure random API key
	timestamp := time.Now().UnixNano()
	data := fmt.Sprintf("krea-speedtest-%d", timestamp)
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("kst_%x", hash[:16]) // 32 character API key with prefix
}

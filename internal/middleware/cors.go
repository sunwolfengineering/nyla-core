package middleware

import (
	"net/http"
	"os"
	"strings"
)

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins   string
	AllowedHeaders   string
	ExposedHeaders   string
	AllowCredentials string
}

// getEnvDefault returns the value of the environment variable or a default
func getEnvDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// NewCORSConfig creates a new CORS configuration from environment variables
func NewCORSConfig() *CORSConfig {
	return &CORSConfig{
		AllowedOrigins:   getEnvDefault("CORS_ALLOWED_ORIGINS", "https://localhost"),
		AllowedHeaders:   getEnvDefault("CORS_ALLOWED_HEADERS", "Content-Type,HX-Request,HX-Target,HX-Current-URL,HX-Trigger,HX-Trigger-Name,HX-History-Restore-Request"),
		ExposedHeaders:   getEnvDefault("CORS_EXPOSED_HEADERS", "HX-Redirect,HX-Location,HX-Push,HX-Refresh,HX-Trigger,HX-Trigger-After-Settle,HX-Trigger-After-Swap"),
		AllowCredentials: getEnvDefault("CORS_ALLOW_CREDENTIALS", "true"),
	}
}

// CORS returns a middleware function that adds CORS headers and handles preflight requests
func (c *CORSConfig) CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		
		// Handle allowed origins
		if c.AllowedOrigins == "*" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		} else if origin != "" {
			// Check if origin is in the allowed origins list
			allowedOriginsList := strings.Split(c.AllowedOrigins, ",")
			for _, allowedOrigin := range allowedOriginsList {
				if strings.TrimSpace(allowedOrigin) == origin {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					w.Header().Set("Vary", "Origin")
					break
				}
			}
		}
		
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", c.AllowedHeaders)
		w.Header().Set("Access-Control-Expose-Headers", c.ExposedHeaders)
		w.Header().Set("Access-Control-Allow-Credentials", c.AllowCredentials)

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

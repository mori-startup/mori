package middleware

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// IPRateLimiter stores rate limiters for different IP addresses
type IPRateLimiter struct {
	ips map[string]*rate.Limiter
	mu  *sync.RWMutex
	r   rate.Limit
	b   int
	lastSeen map[string]time.Time
}

// NewIPRateLimiter creates a new rate limiter for IP addresses
func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	limiter := &IPRateLimiter{
		ips: make(map[string]*rate.Limiter),
		mu:  &sync.RWMutex{},
		r:   r,
		b:   b,
		lastSeen: make(map[string]time.Time),
	}

	// Start cleanup routine
	go limiter.cleanupLoop()

	return limiter
}

// cleanupLoop periodically removes old IP entries
func (i *IPRateLimiter) cleanupLoop() {
	ticker := time.NewTicker(5 * time.Minute)
	for range ticker.C {
		i.mu.Lock()
		now := time.Now()
		for ip, lastSeen := range i.lastSeen {
			// If the IP hasn't been seen in the last 5 minutes, remove it
			if now.Sub(lastSeen) > 5*time.Minute {
				delete(i.ips, ip)
				delete(i.lastSeen, ip)
				log.Printf("Cleaned up rate limiter for IP: %s", ip)
			}
		}
		i.mu.Unlock()
	}
}

// GetLimiter returns the rate limiter for the provided IP address
func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter, exists := i.ips[ip]
	if !exists {
		limiter = rate.NewLimiter(i.r, i.b)
		i.ips[ip] = limiter
	}
	
	// Update last seen time
	i.lastSeen[ip] = time.Now()

	return limiter
}

// getIP extracts the real IP address from the request
func getIP(r *http.Request) string {
	// Try to get IP from X-Forwarded-For header
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		// X-Forwarded-For can contain multiple IPs, take the first one
		ips := strings.Split(ip, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// If no X-Forwarded-For, try X-Real-IP
	ip = r.Header.Get("X-Real-IP")
	if ip != "" {
		return ip
	}

	// If all else fails, use RemoteAddr
	ip = r.RemoteAddr
	// Remove port if present
	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		ip = ip[:idx]
	}
	return ip
}

// shouldSkipRateLimit checks if the request should skip rate limiting
func shouldSkipRateLimit(r *http.Request) bool {
	// Skip rate limiting for OPTIONS requests and static files
	if r.Method == "OPTIONS" || strings.HasPrefix(r.URL.Path, "/imageUpload/") {
		return true
	}

	// Skip rate limiting for specific endpoints
	skipEndpoints := []string{
		"/messageRead",  // Skip rate limiting for message read endpoint
		"/messages",     // Skip rate limiting for message endpoints
		"/unreadMessages",
		"/newMessage",
		"/chatList",
		"/conversationsMsg",
	}

	for _, endpoint := range skipEndpoints {
		if r.URL.Path == endpoint {
			return true
		}
	}

	return false
}

// RateLimit is a middleware that limits requests based on IP address
func RateLimit(next http.Handler) http.Handler {
	// Create rate limiter for normal endpoints
	normalLimiter := NewIPRateLimiter(20, 30)    // 20 req/sec, burst 30

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip rate limiting for certain endpoints
		if shouldSkipRateLimit(r) {
			next.ServeHTTP(w, r)
			return
		}

		// Get the real IP address
		ip := getIP(r)

		// Get the rate limiter for this IP
		ipLimiter := normalLimiter.GetLimiter(ip)

		// Check if the request is allowed
		if !ipLimiter.Allow() {
			log.Printf("Rate limit exceeded for IP: %s on path: %s", ip, r.URL.Path)
			
			// Set CORS headers
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Retry-After", "1")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)

			// Send JSON error response
			errorResponse := map[string]string{
				"error": fmt.Sprintf("Rate limit exceeded. Please wait before trying again. (Path: %s)", r.URL.Path),
				"retry_after": "1",
			}
			json.NewEncoder(w).Encode(errorResponse)
			return
		}

		next.ServeHTTP(w, r)
	})
} 
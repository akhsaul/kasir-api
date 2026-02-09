package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	helper "kasir-api/helpers"
)

// visitor tracks the rate limit state for a single IP.
type visitor struct {
	tokens   float64
	lastSeen time.Time
}

// RateLimiter implements a per-IP token bucket rate limiter.
type RateLimiter struct {
	mu         sync.Mutex
	visitors   map[string]*visitor
	rate       float64 // tokens added per second
	burst      int     // maximum tokens (bucket size)
	cleanupInt time.Duration
}

// NewRateLimiter creates a rate limiter that allows `rate` requests/second with a burst of `burst`.
func NewRateLimiter(rate float64, burst int) *RateLimiter {
	rl := &RateLimiter{
		visitors:   make(map[string]*visitor),
		rate:       rate,
		burst:      burst,
		cleanupInt: 3 * time.Minute,
	}

	// Periodically clean up old visitors
	go rl.cleanup()

	return rl
}

// allow checks if the given IP is allowed to make a request.
func (rl *RateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	now := time.Now()

	if !exists {
		rl.visitors[ip] = &visitor{
			tokens:   float64(rl.burst) - 1, // consume one token
			lastSeen: now,
		}
		return true
	}

	// Add tokens based on elapsed time
	elapsed := now.Sub(v.lastSeen).Seconds()
	v.tokens += elapsed * rl.rate
	if v.tokens > float64(rl.burst) {
		v.tokens = float64(rl.burst)
	}
	v.lastSeen = now

	if v.tokens < 1 {
		return false
	}

	v.tokens--
	return true
}

// cleanup removes stale visitor entries periodically.
func (rl *RateLimiter) cleanup() {
	for {
		time.Sleep(rl.cleanupInt)

		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > rl.cleanupInt {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// extractIP returns the client IP from the request.
func extractIP(r *http.Request) string {
	// Check X-Forwarded-For first (for proxied requests)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

// Limit returns an HTTP middleware that applies rate limiting.
func (rl *RateLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := extractIP(r)
		if !rl.allow(ip) {
			helper.WriteError(w, r, http.StatusTooManyRequests, "Rate limit exceeded. Please try again later.", nil)
			return
		}
		next.ServeHTTP(w, r)
	})
}

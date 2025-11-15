package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/time/rate"
)

// RateLimiter implements token bucket rate limiting
// Security: Prevents abuse and DoS attacks by limiting request rate
type RateLimiter struct {
	visitors map[string]*rate.Limiter
	mu       sync.RWMutex
	rate     rate.Limit
	burst    int
}

// NewRateLimiter creates a new rate limiter
// rate: requests per second
// burst: maximum burst size
func NewRateLimiter(r rate.Limit, b int) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*rate.Limiter),
		rate:     r,
		burst:    b,
	}

	// Cleanup goroutine to remove old visitors
	go rl.cleanupVisitors()

	return rl
}

// getVisitor returns the rate limiter for a given IP
func (rl *RateLimiter) getVisitor(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.visitors[ip]
	if !exists {
		limiter = rate.NewLimiter(rl.rate, rl.burst)
		rl.visitors[ip] = limiter
	}

	return limiter
}

// cleanupVisitors removes inactive visitors periodically
func (rl *RateLimiter) cleanupVisitors() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		// In production, implement proper cleanup based on last access time
		// For now, just clear all to prevent memory leak
		if len(rl.visitors) > 1000 {
			rl.visitors = make(map[string]*rate.Limiter)
		}
		rl.mu.Unlock()
	}
}

// Middleware returns an Echo middleware function
func (rl *RateLimiter) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ip := c.RealIP()
			limiter := rl.getVisitor(ip)

			if !limiter.Allow() {
				return echo.NewHTTPError(http.StatusTooManyRequests, "rate limit exceeded")
			}

			return next(c)
		}
	}
}

// StrictRateLimiter creates a stricter rate limiter for sensitive endpoints
// Example: 5 requests per minute for login endpoints
func StrictRateLimiter() *RateLimiter {
	return NewRateLimiter(rate.Limit(5.0/60.0), 5) // 5 req/min
}

// ModerateRateLimiter creates a moderate rate limiter
// Example: 30 requests per minute
func ModerateRateLimiter() *RateLimiter {
	return NewRateLimiter(rate.Limit(30.0/60.0), 10) // 30 req/min
}

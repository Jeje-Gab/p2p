package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// SecurityHeaders adds common security headers to responses
// Security: Protects against common web vulnerabilities
func SecurityHeaders() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Prevent MIME type sniffing
			c.Response().Header().Set("X-Content-Type-Options", "nosniff")

			// Prevent clickjacking attacks
			c.Response().Header().Set("X-Frame-Options", "DENY")

			// Enable XSS protection (legacy, but doesn't hurt)
			c.Response().Header().Set("X-XSS-Protection", "1; mode=block")

			// Enforce HTTPS (in production with HTTPS enabled)
			// Uncomment when using HTTPS:
			// c.Response().Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

			// Control referrer information
			c.Response().Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

			// Content Security Policy (adjust based on your needs)
			c.Response().Header().Set("Content-Security-Policy", "default-src 'self'")

			return next(c)
		}
	}
}

// CORSConfig returns a CORS configuration for the API
// Security: Properly configured CORS prevents unauthorized cross-origin requests
func CORSConfig(allowOrigins []string) middleware.CORSConfig {
	return middleware.CORSConfig{
		AllowOrigins: allowOrigins,
		AllowMethods: []string{
			echo.GET,
			echo.POST,
			echo.PUT,
			echo.DELETE,
			echo.PATCH,
			echo.OPTIONS,
		},
		AllowHeaders: []string{
			"Authorization",
			"Content-Type",
			"Accept",
		},
		ExposeHeaders: []string{
			"Content-Length",
		},
		AllowCredentials: true,
		MaxAge:           300, // 5 minutes
	}
}

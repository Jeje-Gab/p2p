package middleware

import (
	"net/http"
	"strings"

	"github.com/cs2-p2p-skins/backend/pkg/auth"
	"github.com/labstack/echo/v4"
)

// ContextKey type for context values
type ContextKey string

const (
	// UserIDKey is the context key for user ID
	UserIDKey ContextKey = "user_id"
	// UserEmailKey is the context key for user email
	UserEmailKey ContextKey = "user_email"
	// UserRoleKey is the context key for user role
	UserRoleKey ContextKey = "user_role"
)

// JWTAuth middleware validates JWT tokens and injects user info into context
// Security: Protects endpoints by requiring valid authentication
func JWTAuth(jwtManager *auth.JWTManager) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Extract token from Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing authorization header")
			}

			// Expected format: "Bearer <token>"
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid authorization header format")
			}

			token := parts[1]

			// Validate token
			claims, err := jwtManager.ValidateToken(token)
			if err != nil {
				if err == auth.ErrExpiredToken {
					return echo.NewHTTPError(http.StatusUnauthorized, "token expired")
				}
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
			}

			// Inject user info into context for use by handlers
			c.Set(string(UserIDKey), claims.UserID)
			c.Set(string(UserEmailKey), claims.Email)
			c.Set(string(UserRoleKey), claims.Role)

			return next(c)
		}
	}
}

// RequireRole middleware checks if user has required role
// Security: Implements role-based access control (RBAC)
func RequireRole(roles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userRole := c.Get(string(UserRoleKey))
			if userRole == nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "user role not found in context")
			}

			role := userRole.(string)

			// Check if user has required role
			for _, requiredRole := range roles {
				if role == requiredRole {
					return next(c)
				}
			}

			return echo.NewHTTPError(http.StatusForbidden, "insufficient permissions")
		}
	}
}

// GetUserID extracts user ID from Echo context
func GetUserID(c echo.Context) int64 {
	if userID := c.Get(string(UserIDKey)); userID != nil {
		return userID.(int64)
	}
	return 0
}

// GetUserEmail extracts user email from Echo context
func GetUserEmail(c echo.Context) string {
	if email := c.Get(string(UserEmailKey)); email != nil {
		return email.(string)
	}
	return ""
}

// GetUserRole extracts user role from Echo context
func GetUserRole(c echo.Context) string {
	if role := c.Get(string(UserRoleKey)); role != nil {
		return role.(string)
	}
	return ""
}

package http

import (
	"github.com/cs2-p2p-skins/backend/internal/auth"
	"github.com/cs2-p2p-skins/backend/pkg/middleware"
	pkgAuth "github.com/cs2-p2p-skins/backend/pkg/auth"
	"github.com/labstack/echo/v4"
)

// RegisterRoutes registers all auth routes
func RegisterRoutes(g *echo.Group, usecase auth.UseCase, jwtManager *pkgAuth.JWTManager, rateLimiter *middleware.RateLimiter) {
	handler := NewHandler(usecase)

	// Public routes (no authentication required)
	g.POST("/register", handler.Register, rateLimiter.Middleware())
	g.POST("/login", handler.Login, rateLimiter.Middleware())
	g.POST("/2fa/verify", handler.Verify2FA, rateLimiter.Middleware())
	g.POST("/2fa/verify-backup", handler.VerifyBackupCode, rateLimiter.Middleware())
	g.POST("/2fa/verify-oauth", handler.Verify2FAOAuth, rateLimiter.Middleware())

	// Steam OAuth routes
	g.GET("/steam/login", handler.SteamLogin)
	g.GET("/steam/callback", handler.SteamCallback)

	// Protected routes (require JWT authentication)
	authenticated := g.Group("", middleware.JWTAuth(jwtManager))
	{
		authenticated.GET("/me", handler.GetMe)
		authenticated.POST("/2fa/setup", handler.Setup2FA)
		authenticated.POST("/2fa/enable", handler.Enable2FA)
		authenticated.POST("/2fa/disable", handler.Disable2FA)
	}
}

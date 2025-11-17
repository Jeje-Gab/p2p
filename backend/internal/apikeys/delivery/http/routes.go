package http

import "github.com/labstack/echo/v4"

// RegisterRoutes registers API key management routes
func RegisterRoutes(g *echo.Group, h *Handler) {
	g.POST("", h.CreateAPIKey)
	g.GET("", h.ListAPIKeys)
	g.POST("/:id/revoke", h.RevokeAPIKey)
	g.DELETE("/:id", h.DeleteAPIKey)
}

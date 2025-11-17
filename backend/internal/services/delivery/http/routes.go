package http

import (
	"github.com/labstack/echo/v4"
)

// RegisterRoutes registers service API routes for external systems
func RegisterRoutes(g *echo.Group, h *Handler) {
	// Skins endpoints
	g.GET("/skins/available", h.GetAvailableSkins)
	g.GET("/skins/:id", h.GetSkinByID)

	// Offers endpoints
	g.GET("/offers/open", h.GetOpenOffers)

	// Trades endpoints
	g.GET("/trades/recent", h.GetRecentTrades)

	// Market stats
	g.GET("/market/stats", h.GetMarketStats)
}

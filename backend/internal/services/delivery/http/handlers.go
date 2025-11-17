package http

import (
	"net/http"
	"strconv"

	"github.com/cs2-p2p-skins/backend/internal/offers"
	"github.com/cs2-p2p-skins/backend/internal/skins"
	"github.com/cs2-p2p-skins/backend/internal/trades"
	"github.com/labstack/echo/v4"
)

// Handler handles service API requests for external systems
type Handler struct {
	skinsUC  skins.UseCase
	offersUC offers.UseCase
	tradesUC trades.UseCase
}

// Response is the standard API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// NewHandler creates a new service handler
func NewHandler(skinsUC skins.UseCase, offersUC offers.UseCase, tradesUC trades.UseCase) *Handler {
	return &Handler{
		skinsUC:  skinsUC,
		offersUC: offersUC,
		tradesUC: tradesUC,
	}
}

// GetAvailableSkins returns all available skins for external systems
// GET /services/skins/available
func (h *Handler) GetAvailableSkins(c echo.Context) error {
	ctx := c.Request().Context()

	// Parse pagination
	limit := 100 // Default limit
	offset := 0

	if limitStr := c.QueryParam("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 1000 {
			limit = l
		}
	}

	if offsetStr := c.QueryParam("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	// Get skins with pagination
	skins, err := h.skinsUC.ListSkins(ctx, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "Failed to fetch skins",
		})
	}

	return c.JSON(http.StatusOK, Response{
		Success: true,
		Data: map[string]interface{}{
			"skins":  skins,
			"limit":  limit,
			"offset": offset,
		},
	})
}

// GetSkinByID returns details of a specific skin
// GET /services/skins/:id
func (h *Handler) GetSkinByID(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Invalid skin ID",
		})
	}

	skin, err := h.skinsUC.GetSkin(ctx, id)
	if err != nil {
		return c.JSON(http.StatusNotFound, Response{
			Success: false,
			Message: "Skin not found",
		})
	}

	return c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    skin,
	})
}

// GetOpenOffers returns all open trade offers
// GET /services/offers/open
func (h *Handler) GetOpenOffers(c echo.Context) error {
	ctx := c.Request().Context()

	// Parse pagination
	limit := 100
	offset := 0

	if limitStr := c.QueryParam("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 1000 {
			limit = l
		}
	}

	if offsetStr := c.QueryParam("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	// Get open offers
	offers, err := h.offersUC.ListOffers(ctx, "open", limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "Failed to fetch offers",
		})
	}

	return c.JSON(http.StatusOK, Response{
		Success: true,
		Data: map[string]interface{}{
			"offers": offers,
			"limit":  limit,
			"offset": offset,
		},
	})
}

// GetRecentTrades returns recent completed trades
// GET /services/trades/recent
func (h *Handler) GetRecentTrades(c echo.Context) error {
	ctx := c.Request().Context()

	// Parse pagination
	limit := 50
	offset := 0

	if limitStr := c.QueryParam("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 500 {
			limit = l
		}
	}

	if offsetStr := c.QueryParam("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	// Get all trades
	trades, err := h.tradesUC.ListAllTrades(ctx, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "Failed to fetch trades",
		})
	}

	return c.JSON(http.StatusOK, Response{
		Success: true,
		Data: map[string]interface{}{
			"trades": trades,
			"limit":  limit,
			"offset": offset,
		},
	})
}

// GetMarketStats returns market statistics
// GET /services/market/stats
func (h *Handler) GetMarketStats(c echo.Context) error {
	ctx := c.Request().Context()

	// Get skins for stats (get a large sample)
	skins, err := h.skinsUC.ListSkins(ctx, 10000, 0)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "Failed to fetch market stats",
		})
	}

	// Get trades for volume (get recent trades)
	trades, err := h.tradesUC.ListAllTrades(ctx, 1000, 0)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "Failed to fetch trade stats",
		})
	}

	// Calculate stats
	stats := map[string]interface{}{
		"total_skins":         len(skins),
		"total_trades":        len(trades),
		"available_for_trade": len(skins), // Simplified
	}

	return c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    stats,
	})
}

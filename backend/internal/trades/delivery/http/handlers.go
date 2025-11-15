package http

import (
	"net/http"
	"strconv"

	"github.com/cs2-p2p-skins/backend/internal/trades"
	"github.com/cs2-p2p-skins/backend/pkg/middleware"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	usecase trades.UseCase
}

func NewHandler(uc trades.UseCase) *Handler {
	return &Handler{usecase: uc}
}

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func (h *Handler) GetTrade(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	trade, err := h.usecase.GetTrade(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, APIResponse{
			Success: false,
			Message: "Trade not found",
		})
	}

	return c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    trade,
	})
}

func (h *Handler) ListMyTrades(c echo.Context) error {
	userID := middleware.GetUserID(c)
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	trades, err := h.usecase.ListMyTrades(c.Request().Context(), userID, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Failed to list trades",
		})
	}

	return c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    trades,
	})
}

func (h *Handler) ListAllTrades(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	trades, err := h.usecase.ListAllTrades(c.Request().Context(), limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Failed to list trades",
		})
	}

	return c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    trades,
	})
}

func RegisterRoutes(g *echo.Group, uc trades.UseCase) {
	h := NewHandler(uc)

	g.GET("", h.ListAllTrades)
	g.GET("/me", h.ListMyTrades)
	g.GET("/:id", h.GetTrade)
}

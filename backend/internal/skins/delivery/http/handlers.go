package http

import (
	"net/http"
	"strconv"

	"github.com/cs2-p2p-skins/backend/internal/skins"
	"github.com/cs2-p2p-skins/backend/pkg/middleware"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	usecase skins.UseCase
}

func NewHandler(uc skins.UseCase) *Handler {
	return &Handler{usecase: uc}
}

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// ListSkins returns all available skins
func (h *Handler) ListSkins(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	skins, err := h.usecase.ListSkins(c.Request().Context(), limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Failed to list skins",
		})
	}

	return c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    skins,
	})
}

// GetSkin returns a single skin by ID
func (h *Handler) GetSkin(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Invalid skin ID",
		})
	}

	skin, err := h.usecase.GetSkin(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, APIResponse{
			Success: false,
			Message: "Skin not found",
		})
	}

	return c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    skin,
	})
}

// GetMySkins returns current user's skins
func (h *Handler) GetMySkins(c echo.Context) error {
	userID := middleware.GetUserID(c)

	skins, err := h.usecase.GetUserSkins(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Failed to get user skins",
		})
	}

	return c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    skins,
	})
}

type AddSkinRequest struct {
	SkinID   int64 `json:"skin_id" validate:"required"`
	Quantity int   `json:"quantity" validate:"required,min=1"`
}

// AddSkinToMyInventory adds a skin to user's inventory (for demo purposes)
func (h *Handler) AddSkinToMyInventory(c echo.Context) error {
	userID := middleware.GetUserID(c)

	var req AddSkinRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Invalid request",
		})
	}

	if err := h.usecase.AddSkinToUser(c.Request().Context(), userID, req.SkinID, req.Quantity); err != nil {
		return c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Failed to add skin",
		})
	}

	return c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "Skin added to inventory",
	})
}

func RegisterRoutes(g *echo.Group, uc skins.UseCase) {
	h := NewHandler(uc)

	g.GET("", h.ListSkins)           // GET /skins
	g.GET("/:id", h.GetSkin)         // GET /skins/:id
	g.GET("/me", h.GetMySkins)       // GET /skins/me (requires auth)
	g.POST("/me", h.AddSkinToMyInventory) // POST /skins/me (requires auth)
}

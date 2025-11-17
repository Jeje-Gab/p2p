package http

import (
	"net/http"
	"strconv"

	"github.com/cs2-p2p-skins/backend/internal/apikeys"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	useCase apikeys.UseCase
}

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func NewHandler(useCase apikeys.UseCase) *Handler {
	return &Handler{useCase: useCase}
}

// CreateAPIKey creates a new API key for the authenticated user
// POST /api/api-keys
func (h *Handler) CreateAPIKey(c echo.Context) error {
	ctx := c.Request().Context()

	// Get user ID from JWT context
	userID, ok := c.Get("user_id").(int64)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": "unauthorized",
		})
	}

	// Parse request
	var req apikeys.CreateAPIKeyRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Invalid request body",
		})
	}

	// Validate required fields
	if req.Name == "" {
		return c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "API key name is required",
		})
	}

	// Create API key
	result, err := h.useCase.Create(ctx, userID, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, Response{
		Success: true,
		Message: "API key created successfully. IMPORTANT: Save this key now, it will not be shown again!",
		Data:    result,
	})
}

// ListAPIKeys returns all API keys for the authenticated user
// GET /api/api-keys
func (h *Handler) ListAPIKeys(c echo.Context) error {
	ctx := c.Request().Context()

	// Get user ID from JWT context
	userID, ok := c.Get("user_id").(int64)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": "unauthorized",
		})
	}

	// Get API keys
	keys, err := h.useCase.List(ctx, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    keys,
	})
}

// RevokeAPIKey deactivates an API key
// POST /api/api-keys/:id/revoke
func (h *Handler) RevokeAPIKey(c echo.Context) error {
	ctx := c.Request().Context()

	// Get user ID from JWT context
	userID, ok := c.Get("user_id").(int64)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": "unauthorized",
		})
	}

	// Get key ID from path
	keyID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Invalid API key ID",
		})
	}

	// Revoke API key
	if err := h.useCase.Revoke(ctx, userID, keyID); err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "API key revoked successfully",
	})
}

// DeleteAPIKey permanently removes an API key
// DELETE /api/api-keys/:id
func (h *Handler) DeleteAPIKey(c echo.Context) error {
	ctx := c.Request().Context()

	// Get user ID from JWT context
	userID, ok := c.Get("user_id").(int64)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": "unauthorized",
		})
	}

	// Get key ID from path
	keyID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Invalid API key ID",
		})
	}

	// Delete API key
	if err := h.useCase.Delete(ctx, userID, keyID); err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "API key deleted successfully",
	})
}

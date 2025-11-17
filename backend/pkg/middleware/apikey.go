package middleware

import (
	"net/http"

	"github.com/cs2-p2p-skins/backend/internal/apikeys"
	"github.com/labstack/echo/v4"
)

// APIKeyAuth middleware validates API keys for external systems using database
// Security: Allows external systems to access public service endpoints
func APIKeyAuth(apiKeyUseCase apikeys.UseCase) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Extract API key from header
			apiKey := c.Request().Header.Get("X-API-Key")
			if apiKey == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, map[string]interface{}{
					"success": false,
					"message": "missing API key - provide X-API-Key header",
				})
			}

			// Validate API key against database
			validatedKey, err := apiKeyUseCase.Validate(c.Request().Context(), apiKey)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, map[string]interface{}{
					"success": false,
					"message": err.Error(),
				})
			}

			// API key is valid, inject metadata into context
			c.Set("api_key_auth", true)
			c.Set("api_key_id", validatedKey.ID)
			c.Set("api_key_user_id", validatedKey.UserID)
			c.Set("auth_type", "api_key")

			return next(c)
		}
	}
}

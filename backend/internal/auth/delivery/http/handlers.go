package http

import (
	"net/http"

	"github.com/cs2-p2p-skins/backend/internal/auth"
	"github.com/cs2-p2p-skins/backend/pkg/middleware"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	usecase auth.UseCase
}

func NewHandler(usecase auth.UseCase) *Handler {
	return &Handler{usecase: usecase}
}

// RegisterRequest represents the registration request
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// LoginRequest represents the login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// Verify2FARequest represents the 2FA verification request
type Verify2FARequest struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required,len=6"`
}

// Enable2FARequest represents the request to enable 2FA
type Enable2FARequest struct {
	Code string `json:"code" validate:"required,len=6"`
}

// APIResponse is a generic API response wrapper
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

// Register handles user registration
// POST /auth/register
func (h *Handler) Register(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Invalid request",
			Errors:  err.Error(),
		})
	}

	user, err := h.usecase.Register(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		statusCode := http.StatusInternalServerError
		message := "Failed to register user"

		// Handle specific errors
		if err.Error() == "user already exists" {
			statusCode = http.StatusConflict
			message = "User already exists"
		}

		return c.JSON(statusCode, APIResponse{
			Success: false,
			Message: message,
			Errors:  err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Message: "User registered successfully",
		Data: map[string]interface{}{
			"id":    user.ID,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

// Login handles user login (first factor)
// POST /auth/login
func (h *Handler) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Invalid request",
			Errors:  err.Error(),
		})
	}

	ip := c.RealIP()

	requiresTwoFA, token, user, err := h.usecase.Login(c.Request().Context(), req.Email, req.Password, ip)
	if err != nil {
		statusCode := http.StatusUnauthorized
		message := "Invalid credentials"

		if err.Error() == "account temporarily locked due to too many failed attempts" {
			statusCode = http.StatusTooManyRequests
			message = "Account locked due to too many failed login attempts. Please try again later."
		}

		return c.JSON(statusCode, APIResponse{
			Success: false,
			Message: message,
		})
	}

	// If 2FA is required, return special response
	if requiresTwoFA {
		return c.JSON(http.StatusOK, APIResponse{
			Success: true,
			Message: "2FA code required",
			Data: map[string]interface{}{
				"requires_2fa": true,
			},
		})
	}

	return c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "Login successful",
		Data: map[string]interface{}{
			"token": token,
			"user": map[string]interface{}{
				"id":            user.ID,
				"email":         user.Email,
				"steam_id":      user.SteamID,
				"role":          user.Role,
				"twofa_enabled": user.TwoFAEnabled,
				"created_at":    user.CreatedAt,
				"updated_at":    user.UpdatedAt,
			},
		},
	})
}

// Verify2FA handles 2FA verification
// POST /auth/2fa/verify
func (h *Handler) Verify2FA(c echo.Context) error {
	var req Verify2FARequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Invalid request",
			Errors:  err.Error(),
		})
	}

	ip := c.RealIP()

	token, user, err := h.usecase.Verify2FA(c.Request().Context(), req.Email, req.Code, ip)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, APIResponse{
			Success: false,
			Message: "Invalid 2FA code",
		})
	}

	return c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "2FA verification successful",
		Data: map[string]interface{}{
			"token": token,
			"user": map[string]interface{}{
				"id":            user.ID,
				"email":         user.Email,
				"steam_id":      user.SteamID,
				"role":          user.Role,
				"twofa_enabled": user.TwoFAEnabled,
				"created_at":    user.CreatedAt,
				"updated_at":    user.UpdatedAt,
			},
		},
	})
}

// GetMe returns the current authenticated user
// GET /auth/me
func (h *Handler) GetMe(c echo.Context) error {
	userID := middleware.GetUserID(c)

	user, err := h.usecase.GetCurrentUser(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, APIResponse{
			Success: false,
			Message: "User not found",
		})
	}

	return c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"id":            user.ID,
			"email":         user.Email,
			"steam_id":      user.SteamID,
			"role":          user.Role,
			"twofa_enabled": user.TwoFAEnabled,
			"created_at":    user.CreatedAt,
		},
	})
}

// Setup2FA generates a new 2FA secret
// POST /auth/2fa/setup
func (h *Handler) Setup2FA(c echo.Context) error {
	userID := middleware.GetUserID(c)

	secret, qrURL, err := h.usecase.Setup2FA(c.Request().Context(), userID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		message := "Failed to setup 2FA"

		if err.Error() == "2FA already enabled" {
			statusCode = http.StatusConflict
			message = "2FA is already enabled"
		}

		return c.JSON(statusCode, APIResponse{
			Success: false,
			Message: message,
		})
	}

	return c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "2FA setup initiated. Scan QR code with Google Authenticator and verify.",
		Data: map[string]interface{}{
			"secret": secret,
			"qr_url": qrURL,
		},
	})
}

// Enable2FA enables 2FA after code verification
// POST /auth/2fa/enable
func (h *Handler) Enable2FA(c echo.Context) error {
	userID := middleware.GetUserID(c)

	var req Enable2FARequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Invalid request",
			Errors:  err.Error(),
		})
	}

	if err := h.usecase.Enable2FA(c.Request().Context(), userID, req.Code); err != nil {
		return c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Invalid 2FA code or 2FA already enabled",
		})
	}

	return c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "2FA enabled successfully",
	})
}

// Disable2FA disables 2FA after code verification
// POST /auth/2fa/disable
func (h *Handler) Disable2FA(c echo.Context) error {
	userID := middleware.GetUserID(c)

	var req Enable2FARequest // Reuse same struct
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Invalid request",
			Errors:  err.Error(),
		})
	}

	if err := h.usecase.Disable2FA(c.Request().Context(), userID, req.Code); err != nil {
		return c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Invalid 2FA code or 2FA not enabled",
		})
	}

	return c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "2FA disabled successfully",
	})
}

// SteamLogin redirects to Steam OAuth
// GET /auth/steam/login
func (h *Handler) SteamLogin(c echo.Context) error {
	authURL, state, err := h.usecase.GenerateSteamAuthURL(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Failed to generate Steam login URL",
		})
	}

	// Store state in session/cookie for CSRF validation
	c.SetCookie(&http.Cookie{
		Name:     "steam_oauth_state",
		Value:    state,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
		MaxAge:   300, // 5 minutes
	})

	return c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "Redirect to Steam for authentication",
		Data: map[string]interface{}{
			"auth_url": authURL,
		},
	})
}

// SteamCallback handles Steam OAuth callback
// GET /auth/steam/callback
func (h *Handler) SteamCallback(c echo.Context) error {
	// Extract all query parameters
	values := make(map[string]string)
	for key, vals := range c.QueryParams() {
		if len(vals) > 0 {
			values[key] = vals[0]
		}
	}

	token, err := h.usecase.HandleSteamCallback(c.Request().Context(), values)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, APIResponse{
			Success: false,
			Message: "Steam authentication failed",
			Errors:  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "Steam login successful",
		Data: map[string]interface{}{
			"token": token,
		},
	})
}

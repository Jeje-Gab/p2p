package http

import (
	"log"
	"net/http"
	"strconv"

	"github.com/cs2-p2p-skins/backend/internal/offers"
	"github.com/cs2-p2p-skins/backend/pkg/middleware"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	usecase offers.UseCase
}

func NewHandler(uc offers.UseCase) *Handler {
	return &Handler{usecase: uc}
}

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

type CreateOfferRequest struct {
	SkinOfferedID   int64 `json:"skin_offered_id" validate:"required"`
	SkinRequestedID int64 `json:"skin_requested_id" validate:"required"`
}

func (h *Handler) CreateOffer(c echo.Context) error {
	userID := middleware.GetUserID(c)

	var req CreateOfferRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Invalid request",
		})
	}

	offer, err := h.usecase.CreateOffer(c.Request().Context(), userID, req.SkinOfferedID, req.SkinRequestedID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Message: "Offer created successfully",
		Data:    offer,
	})
}

func (h *Handler) GetOffer(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	offer, err := h.usecase.GetOffer(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, APIResponse{
			Success: false,
			Message: "Offer not found",
		})
	}

	return c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    offer,
	})
}

func (h *Handler) ListOffers(c echo.Context) error {
	status := c.QueryParam("status")
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	offers, err := h.usecase.ListOffers(c.Request().Context(), status, limit, offset)
	if err != nil {
		log.Printf("ERROR ListOffers: status=%s, limit=%d, offset=%d, error=%v", status, limit, offset, err)
		return c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Failed to list offers",
			Errors:  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    offers,
	})
}

func (h *Handler) ListMyOffers(c echo.Context) error {
	userID := middleware.GetUserID(c)
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	offers, err := h.usecase.ListMyOffers(c.Request().Context(), userID, limit, offset)
	if err != nil {
		log.Printf("ERROR ListMyOffers: userID=%d, limit=%d, offset=%d, error=%v", userID, limit, offset, err)
		return c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Failed to list offers",
			Errors:  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    offers,
	})
}

func (h *Handler) AcceptOffer(c echo.Context) error {
	userID := middleware.GetUserID(c)
	offerID, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	if err := h.usecase.AcceptOffer(c.Request().Context(), offerID, userID); err != nil {
		return c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "Offer accepted successfully",
	})
}

func (h *Handler) CancelOffer(c echo.Context) error {
	userID := middleware.GetUserID(c)
	offerID, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	if err := h.usecase.CancelOffer(c.Request().Context(), offerID, userID); err != nil {
		return c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "Offer cancelled successfully",
	})
}

func RegisterRoutes(g *echo.Group, uc offers.UseCase) {
	h := NewHandler(uc)

	g.POST("", h.CreateOffer)
	g.GET("", h.ListOffers)
	g.GET("/my", h.ListMyOffers)
	g.GET("/:id", h.GetOffer)
	g.POST("/:id/accept", h.AcceptOffer)
	g.POST("/:id/cancel", h.CancelOffer)
}

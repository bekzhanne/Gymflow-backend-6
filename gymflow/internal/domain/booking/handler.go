package booking

import (
	"net/http"

	"gymflow/internal/middleware"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// POST /api/v1/classes (admin/trainer)
func (h *Handler) CreateClass(c *gin.Context) {
	var req CreateClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	class, err := h.service.CreateClass(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, ToClassResponse(class))
}

// GET /api/v1/classes
func (h *Handler) ListClasses(c *gin.Context) {
	classes, err := h.service.ListClasses()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list classes"})
		return
	}
	resp := make([]*ClassResponse, 0, len(classes))
	for i := range classes {
		resp = append(resp, ToClassResponse(&classes[i]))
	}
	c.JSON(http.StatusOK, resp)
}

// POST /api/v1/bookings
func (h *Handler) CreateBooking(c *gin.Context) {
	var req CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userIDAny, _ := c.Get(middleware.ContextUserIDKey)
	userID := userIDAny.(uint)

	b, err := h.service.CreateBooking(userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, ToBookingResponse(b))
}

// GET /api/v1/bookings
func (h *Handler) ListBookings(c *gin.Context) {
	userIDAny, _ := c.Get(middleware.ContextUserIDKey)
	userID := userIDAny.(uint)

	bookings, err := h.service.ListBookings(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list bookings"})
		return
	}
	resp := make([]*BookingResponse, 0, len(bookings))
	for i := range bookings {
		resp = append(resp, ToBookingResponse(&bookings[i]))
	}
	c.JSON(http.StatusOK, resp)
}

// POST /api/v1/bookings/:id/cancel
func (h *Handler) CancelBooking(c *gin.Context) {
	var uri struct {
		ID uint `uri:"id" binding:"required"`
	}
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userIDAny, _ := c.Get(middleware.ContextUserIDKey)
	userID := userIDAny.(uint)

	b, err := h.service.CancelBooking(userID, uri.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ToBookingResponse(b))
}

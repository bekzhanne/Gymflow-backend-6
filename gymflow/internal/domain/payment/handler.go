package payment

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

// POST /api/v1/payments
func (h *Handler) CreatePayment(c *gin.Context) {
	var req CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userIDAny, _ := c.Get(middleware.ContextUserIDKey)
	userID := userIDAny.(uint)

	p, err := h.service.CreatePayment(userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, ToPaymentResponse(p))
}

// GET /api/v1/payments
func (h *Handler) ListPayments(c *gin.Context) {
	userIDAny, _ := c.Get(middleware.ContextUserIDKey)
	userID := userIDAny.(uint)

	payments, err := h.service.ListPayments(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list payments"})
		return
	}

	resp := make([]*PaymentResponse, 0, len(payments))
	for i := range payments {
		resp = append(resp, ToPaymentResponse(&payments[i]))
	}
	c.JSON(http.StatusOK, resp)
}

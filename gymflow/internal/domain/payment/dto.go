package payment

type CreatePaymentRequest struct {
	BookingID uint    `json:"booking_id" binding:"required"`
	Amount    float64 `json:"amount" binding:"required,gt=0"`
	Method    string  `json:"method" binding:"required"` // card, cash, etc.
}

type PaymentResponse struct {
	ID        uint    `json:"id"`
	UserID    uint    `json:"user_id"`
	BookingID uint    `json:"booking_id"`
	Amount    float64 `json:"amount"`
	Status    string  `json:"status"`
	Method    string  `json:"method"`
}

func ToPaymentResponse(p *Payment) *PaymentResponse {
	return &PaymentResponse{
		ID:        p.ID,
		UserID:    p.UserID,
		BookingID: p.BookingID,
		Amount:    p.Amount,
		Status:    p.Status,
		Method:    p.Method,
	}
}

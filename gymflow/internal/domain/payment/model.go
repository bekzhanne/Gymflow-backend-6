package payment

import "time"

type Payment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uint      `json:"user_id"`
	BookingID uint      `json:"booking_id"`
	Amount    float64   `json:"amount"`
	Status    string    `json:"status"`
	Method    string    `json:"method"`
}

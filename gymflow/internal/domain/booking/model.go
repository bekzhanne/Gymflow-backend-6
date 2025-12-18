package booking

import (
	"time"
)

const (
	BookingStatusBooked    = "booked"
	BookingStatusCancelled = "cancelled"
	BookingStatusWaitlist  = "waitlist"

	PaymentStatusPending = "pending"
	PaymentStatusPaid    = "paid"
)

type GymClass struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	TrainerID   uint      `json:"trainer_id"`
	Capacity    int       `json:"capacity"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Price       float64   `json:"price"`
}

type Booking struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	UserID        uint      `json:"user_id"`
	ClassID       uint      `json:"class_id"`
	Status        string    `json:"status"`
	PaymentStatus string    `json:"payment_status"`
}

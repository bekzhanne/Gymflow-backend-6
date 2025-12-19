package booking

type CreateClassRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	TrainerID   uint    `json:"trainer_id" binding:"required"`
	Capacity    int     `json:"capacity" binding:"required,min=1"`
	StartTime   string  `json:"start_time" binding:"required"` // ISO8601
	EndTime     string  `json:"end_time" binding:"required"`
	Price       float64 `json:"price" binding:"required,min=0"`
}

type ClassResponse struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	TrainerID   uint    `json:"trainer_id"`
	Capacity    int     `json:"capacity"`
	StartTime   string  `json:"start_time"`
	EndTime     string  `json:"end_time"`
	Price       float64 `json:"price"`
}

type CreateBookingRequest struct {
	ClassID uint `json:"class_id" binding:"required"`
}

type BookingResponse struct {
	ID            uint   `json:"id"`
	UserID        uint   `json:"user_id"`
	ClassID       uint   `json:"class_id"`
	Status        string `json:"status"`
	PaymentStatus string `json:"payment_status"`
}

func ToClassResponse(c *GymClass) *ClassResponse {
	return &ClassResponse{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
		TrainerID:   c.TrainerID,
		Capacity:    c.Capacity,
		StartTime:   c.StartTime.Format("2006-01-02T15:04:05Z07:00"),
		EndTime:     c.EndTime.Format("2006-01-02T15:04:05Z07:00"),
		Price:       c.Price,
	}
}

func ToBookingResponse(b *Booking) *BookingResponse {
	return &BookingResponse{
		ID:            b.ID,
		UserID:        b.UserID,
		ClassID:       b.ClassID,
		Status:        b.Status,
		PaymentStatus: b.PaymentStatus,
	}
}
	
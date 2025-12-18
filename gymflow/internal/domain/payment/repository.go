package payment

import "gorm.io/gorm"

type Repository interface {
	Create(p *Payment) error
	ListByUser(userID uint) ([]Payment, error)
	FindByBookingID(bookingID uint) (*Payment, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(p *Payment) error {
	return r.db.Create(p).Error
}

func (r *repository) FindByBookingID(bookingID uint) (*Payment, error) {
	var p Payment
	err := r.db.Where("booking_id = ?", bookingID).First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *repository) ListByUser(userID uint) ([]Payment, error) {
	var pay []Payment
	if err := r.db.Where("user_id = ?", userID).Find(&pay).Error; err != nil {
		return nil, err
	}
	return pay, nil
}

package booking

import "gorm.io/gorm"

type Repository interface {
	CreateClass(c *GymClass) error
	ListClasses() ([]GymClass, error)
	FindClassByID(id uint) (*GymClass, error)

	CreateBooking(b *Booking) error
	ListBookingsByUser(userID uint) ([]Booking, error)
	CountBookingsForClass(classID uint) (int64, error)
	UpdateBooking(b *Booking) error
	FindBookingByID(id uint) (*Booking, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateClass(c *GymClass) error {
	return r.db.Create(c).Error
}

func (r *repository) ListClasses() ([]GymClass, error) {
	var classes []GymClass
	if err := r.db.Find(&classes).Error; err != nil {
		return nil, err
	}
	return classes, nil
}

func (r *repository) FindClassByID(id uint) (*GymClass, error) {
	var c GymClass
	if err := r.db.First(&c, id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *repository) CreateBooking(b *Booking) error {
	return r.db.Create(b).Error
}

func (r *repository) ListBookingsByUser(userID uint) ([]Booking, error) {
	var bookings []Booking
	if err := r.db.Where("user_id = ?", userID).Find(&bookings).Error; err != nil {
		return nil, err
	}
	return bookings, nil
}

func (r *repository) CountBookingsForClass(classID uint) (int64, error) {
	var count int64
	err := r.db.Model(&Booking{}).
		Where("class_id = ? AND status = ?", classID, BookingStatusBooked).
		Count(&count).Error
	return count, err
}

func (r *repository) UpdateBooking(b *Booking) error {
	return r.db.Save(b).Error
}

func (r *repository) FindBookingByID(id uint) (*Booking, error) {
	var b Booking
	if err := r.db.First(&b, id).Error; err != nil {
		return nil, err
	}
	return &b, nil
}

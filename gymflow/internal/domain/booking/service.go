package booking

import (
	"errors"
	"time"
)

type Service interface {
	CreateClass(req CreateClassRequest) (*GymClass, error)
	ListClasses() ([]GymClass, error)
	CreateBooking(userID uint, req CreateBookingRequest) (*Booking, error)
	ListBookings(userID uint) ([]Booking, error)
	CancelBooking(userID, bookingID uint) (*Booking, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateClass(req CreateClassRequest) (*GymClass, error) {
	start, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		return nil, err
	}
	end, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		return nil, err
	}
	c := &GymClass{
		Name:        req.Name,
		Description: req.Description,
		TrainerID:   req.TrainerID,
		Capacity:    req.Capacity,
		StartTime:   start,
		EndTime:     end,
		Price:       req.Price,
	}
	if err := s.repo.CreateClass(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *service) ListClasses() ([]GymClass, error) {
	return s.repo.ListClasses()
}

func (s *service) CreateBooking(userID uint, req CreateBookingRequest) (*Booking, error) {
	class, err := s.repo.FindClassByID(req.ClassID)
	if err != nil {
		return nil, err
	}

	count, err := s.repo.CountBookingsForClass(class.ID)
	if err != nil {
		return nil, err
	}

	status := BookingStatusBooked
	if int(count) >= class.Capacity {
		status = BookingStatusWaitlist
	}

	b := &Booking{
		UserID:        userID,
		ClassID:       class.ID,
		Status:        status,
		PaymentStatus: PaymentStatusPending,
	}
	if err := s.repo.CreateBooking(b); err != nil {
		return nil, err
	}
	return b, nil
}

func (s *service) ListBookings(userID uint) ([]Booking, error) {
	return s.repo.ListBookingsByUser(userID)
}

func (s *service) CancelBooking(userID, bookingID uint) (*Booking, error) {
	b, err := s.repo.FindBookingByID(bookingID)
	if err != nil {
		return nil, err
	}
	if b.UserID != userID {
		return nil, errors.New("cannot cancel foreign booking")
	}
	b.Status = BookingStatusCancelled
	if err := s.repo.UpdateBooking(b); err != nil {
		return nil, err
	}
	return b, nil
}

package admin

import (
	"time"

	"gymflow/internal/domain/booking"
	"gymflow/internal/domain/payment"
	"gymflow/internal/domain/user"
	"gorm.io/gorm"
)

type Service interface {
	GetDashboard() (*DashboardResponse, error)
}

type service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) Service {
	return &service{db: db}
}

func (s *service) GetDashboard() (*DashboardResponse, error) {
	var resp DashboardResponse

	s.db.Model(&user.User{}).Count(&resp.TotalUsers)
	s.db.Model(&booking.GymClass{}).Count(&resp.TotalClasses)
	s.db.Model(&booking.Booking{}).Count(&resp.TotalBookings)

	s.db.Model(&user.User{}).Where("active = ?", true).Count(&resp.ActiveMembers)
	s.db.Model(&booking.GymClass{}).Where("start_time > ?", time.Now()).Count(&resp.UpcomingClasses)

	type res struct {
		Sum float64
	}
	var r res
	s.db.Model(&payment.Payment{}).Select("COALESCE(sum(amount),0) as sum").Scan(&r)
	resp.TotalRevenue = r.Sum

	return &resp, nil
}

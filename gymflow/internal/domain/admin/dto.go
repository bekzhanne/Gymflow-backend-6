package admin

type DashboardResponse struct {
	TotalUsers     int64   `json:"total_users"`
	TotalClasses   int64   `json:"total_classes"`
	TotalBookings  int64   `json:"total_bookings"`
	TotalRevenue   float64 `json:"total_revenue"`
	ActiveMembers  int64   `json:"active_members"`
	UpcomingClasses int64  `json:"upcoming_classes"`
}

package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gymflow/internal/config"
	"gymflow/internal/domain/admin"
	"gymflow/internal/domain/auth"
	"gymflow/internal/domain/booking"
	"gymflow/internal/domain/payment"
	"gymflow/internal/domain/user"
	"gymflow/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// Auto migrate all models
	db.AutoMigrate(
		&user.User{},
		&booking.GymClass{},
		&booking.Booking{},
		&payment.Payment{},
	)

	return db
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	db := setupTestDB()

	// Config
	cfg := &config.Config{
		JWTSecret:   "test-secret-key",
		JWTTTLHours: 72,
	}

	// Repositories
	userRepo := user.NewRepository(db)
	bookingRepo := booking.NewRepository(db)
	paymentRepo := payment.NewRepository(db)

	// Services
	userService := user.NewService(userRepo)
	bookingService := booking.NewService(bookingRepo)
	paymentService := payment.NewService(paymentRepo)
	adminService := admin.NewService(db)

	// Handlers
	userHandler := user.NewHandler(cfg, userService)
	authHandler := auth.NewHandler(userHandler)
	bookingHandler := booking.NewHandler(bookingService)
	paymentHandler := payment.NewHandler(paymentService)
	adminHandler := admin.NewHandler(adminService)

	// Router
	r := gin.New()

	// Public routes
	api := r.Group("/api/v1")
	{
		// Auth
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/register", authHandler.Register)
			authGroup.POST("/login", authHandler.Login)
		}

		// Public: list classes
		api.GET("/classes", bookingHandler.ListClasses)
	}

	// Protected routes (any authenticated user)
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(cfg))
	{
		// User routes
		protected.GET("/users", userHandler.ListUsers)

		// Booking routes
		protected.POST("/bookings", bookingHandler.CreateBooking)
		protected.GET("/bookings", bookingHandler.ListBookings)
		protected.POST("/bookings/:id/cancel", bookingHandler.CancelBooking)

		// Payment routes
		protected.POST("/payments", paymentHandler.CreatePayment)
		protected.GET("/payments", paymentHandler.ListPayments)
	}

	// Trainer/Admin routes
	trainerRoutes := api.Group("")
	trainerRoutes.Use(middleware.AuthMiddleware(cfg, "trainer", "admin"))
	{
		trainerRoutes.POST("/classes", bookingHandler.CreateClass)
	}

	// Admin only routes
	adminRoutes := api.Group("/admin")
	adminRoutes.Use(middleware.AuthMiddleware(cfg, "admin"))
	{
		adminRoutes.GET("/dashboard", adminHandler.Dashboard)
	}

	return r
}

func makeRequest(t *testing.T, router *gin.Engine, method, url string, body interface{}, token string) *httptest.ResponseRecorder {
	var reqBody *bytes.Buffer

	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("Failed to marshal request body: %v", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	} else {
		reqBody = bytes.NewBuffer([]byte{})
	}

	req, err := http.NewRequest(method, url, reqBody)
	assert.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	return w
}

// seedTestData - helper для создания тестовых данных
func seedTestData(db *gorm.DB) {
	// Admin user
	adminUser := &user.User{
		Name:           "Admin User",
		Email:          "admin@example.com",
		PasswordHash:   "$2a$10$...", // bcrypt hash для "password"
		Role:           user.RoleAdmin,
		MembershipTier: user.MembershipBasic,
		Active:         true,
	}
	db.Create(adminUser)

	// Member user
	memberUser := &user.User{
		Name:           "Member User",
		Email:          "member@example.com",
		PasswordHash:   "$2a$10$...", // bcrypt hash для "password"
		Role:           user.RoleMember,
		MembershipTier: user.MembershipBasic,
		Active:         true,
	}
	db.Create(memberUser)

	// Gym Class
	gymClass := &booking.GymClass{
		Name:        "Test Yoga",
		Description: "Beginner yoga class",
		TrainerID:   1,
		Capacity:    10,
		Price:       50.0,
	}
	db.Create(gymClass)

	// Booking
	bookingRecord := &booking.Booking{
		UserID:        memberUser.ID,
		ClassID:       gymClass.ID,
		Status:        booking.BookingStatusBooked,
		PaymentStatus: booking.PaymentStatusPaid,
	}
	db.Create(bookingRecord)

	// Payment
	paymentRecord := &payment.Payment{
		UserID:    memberUser.ID,
		BookingID: bookingRecord.ID,
		Amount:    50.0,
		Status:    "paid",
		Method:    "card",
	}
	db.Create(paymentRecord)
}
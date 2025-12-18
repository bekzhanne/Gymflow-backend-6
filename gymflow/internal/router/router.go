package router

import (
	"gymflow/internal/config"
	"gymflow/internal/database"
	"gymflow/internal/domain/admin"
	"gymflow/internal/domain/booking"
	"gymflow/internal/domain/payment"
	"gymflow/internal/domain/user"
	"gymflow/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(cfg *config.Config, db *gorm.DB) *gin.Engine {
	r := gin.New()
	r.Use(middleware.Logger(), gin.Recovery())

	// Repos & services
	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(cfg, userService)

	bookingRepo := booking.NewRepository(db)
	bookingService := booking.NewService(bookingRepo)
	bookingHandler := booking.NewHandler(bookingService)

	paymentRepo := payment.NewRepository(db)
	paymentService := payment.NewService(paymentRepo)
	paymentHandler := payment.NewHandler(paymentService)

	adminService := admin.NewService(db)
	adminHandler := admin.NewHandler(adminService)

	api := r.Group("/api/v1")

	

	// Auth
	api.POST("/auth/register", userHandler.Register)
	api.POST("/auth/login", userHandler.Login)

	// Public classes
	api.GET("/classes", bookingHandler.ListClasses)

	// Authenticated routes
	authMember := api.Group("/")
	authMember.Use(middleware.AuthMiddleware(cfg, user.RoleMember, user.RoleTrainer, user.RoleAdmin))

	authMember.POST("/bookings", bookingHandler.CreateBooking)
	authMember.GET("/bookings", bookingHandler.ListBookings)
	authMember.POST("/bookings/:id/cancel", bookingHandler.CancelBooking)

	authMember.POST("/payments", paymentHandler.CreatePayment)
	authMember.GET("/payments", paymentHandler.ListPayments)

	// Trainer/Admin
	authTrainer := api.Group("/")
	authTrainer.Use(middleware.AuthMiddleware(cfg, user.RoleTrainer, user.RoleAdmin))
	authTrainer.POST("/classes", bookingHandler.CreateClass)

	// Admin only
	authAdmin := api.Group("/admin")
	authAdmin.Use(middleware.AuthMiddleware(cfg, user.RoleAdmin))
	authAdmin.GET("/dashboard", adminHandler.Dashboard)

	// Healthcheck
	r.GET("/health", func(c *gin.Context) {
		if err := database.PingRedis(c, database.NewRedisClient(cfg)); err != nil {
			c.JSON(500, gin.H{"status": "redis down", "error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"status": "ok"})
	})

	return r
}

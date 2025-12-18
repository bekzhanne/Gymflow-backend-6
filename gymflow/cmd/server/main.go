package main

import (
	"log"

	"gymflow/internal/config"
	"gymflow/internal/database"
	"gymflow/internal/domain/booking"
	"gymflow/internal/domain/payment"
	"gymflow/internal/domain/user"
	"gymflow/internal/router"
)

func main() {
	cfg := config.LoadConfig()

	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("cannot connect db: %v", err)
	}

	// миграция всех моделей
	if err := db.AutoMigrate(
		&user.User{},
		&booking.GymClass{},
		&booking.Booking{},
		&payment.Payment{},
	); err != nil {
		log.Fatalf("auto migrate failed: %v", err)
	}

	r := router.SetupRouter(cfg, db)

	log.Printf("GymFlow running on :%s", cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
	
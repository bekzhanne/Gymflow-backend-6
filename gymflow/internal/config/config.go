package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	RedisAddr     string
	RedisPassword string
	RedisDB       int

	JWTSecret   string
	JWTTTLHours int
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		AppPort:      getEnv("APP_PORT", "8080"),
		DBHost:       getEnv("DB_HOST", "localhost"),
		DBPort:       getEnv("DB_PORT", "5432"),
		DBUser:       getEnv("DB_USER", "gymflow"),
		DBPassword:   getEnv("DB_PASSWORD", "gymflow"),
		DBName:       getEnv("DB_NAME", "gymflow"),
		RedisAddr:    getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		JWTSecret:    getEnv("JWT_SECRET", "changeme"),
	}

	redisDB, err := strconv.Atoi(getEnv("REDIS_DB", "0"))
	if err != nil {
		log.Fatalf("invalid REDIS_DB: %v", err)
	}
	cfg.RedisDB = redisDB

	ttl, err := strconv.Atoi(getEnv("JWT_TTL_HOURS", "24"))
	if err != nil {
		log.Fatalf("invalid JWT_TTL_HOURS: %v", err)
	}
	cfg.JWTTTLHours = ttl

	return cfg
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

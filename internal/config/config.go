package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port       string
	MongoURI   string
	DBName     string
	JWTSecret  string
}

var AppConfig *Config

func Load() error {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	AppConfig = &Config{
		Port:      getEnv("PORT", "8080"),
		MongoURI:  getEnv("MONGODB_URI", ""),
		DBName:    getEnv("DB_NAME", "onestay"),
		JWTSecret: getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
	}

	if AppConfig.MongoURI == "" {
		log.Fatal("MONGODB_URI is required")
	}

	if AppConfig.JWTSecret == "your-secret-key-change-in-production" {
		log.Println("Warning: Using default JWT_SECRET. Change it in production!")
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

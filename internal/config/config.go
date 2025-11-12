package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	Environment string
	MongoURI    string
	DBName      string
	JWTSecret   string

	// Razorpay
	RazorpayKey    string
	RazorpaySecret string

	// Agora
	AgoraAppID   string
	AgoraAppCert string

	// Email
	SMTPHost     string
	SMTPPort     string
	SMTPUser     string
	SMTPPassword string
}

func LoadConfig() *Config {
	// Load .env file if it exists
	_ = godotenv.Load()

	// Set default values
	config := &Config{
		Port:           getEnv("PORT", "8080"),
		Environment:    getEnv("ENVIRONMENT", "development"),
		MongoURI:       getEnv("MONGO_URI", "mongodb://localhost:27017"),
		DBName:         getEnv("DB_NAME", "hsb_backend"),
		JWTSecret:      getEnv("JWT_SECRET", "default-secret-key"),
		RazorpayKey:    getEnv("RAZORPAY_KEY", ""),
		RazorpaySecret: getEnv("RAZORPAY_SECRET", ""),
		AgoraAppID:     getEnv("AGORA_APP_ID", ""),
		AgoraAppCert:   getEnv("AGORA_APP_CERT", ""),
		SMTPHost:       getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:       getEnv("SMTP_PORT", "587"),
		SMTPUser:       getEnv("SMTP_USER", ""),
		SMTPPassword:   getEnv("SMTP_PASSWORD", ""),
	}

	// Validate required configurations
	if config.JWTSecret == "default-secret-key" && config.Environment == "production" {
		log.Fatal("JWT_SECRET must be set in production environment")
	}

	return config
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}

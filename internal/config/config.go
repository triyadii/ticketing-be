package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName        string
	AppPort        string
	AppEnv         string
	AppURL         string
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	DBSSLMode      string
	DBTimeZone     string
	JWTSecret      string
	JWTExpireHours int
	UploadDir      string
	MaxFileSizeMB  int
}

var AppConfig *Config

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Note: .env file not found, reading system environment variables")
	}

	jwtExpire, err := strconv.Atoi(getEnv("JWT_EXPIRE_HOURS", "72"))
	if err != nil {
		jwtExpire = 72
	}

	maxFileSize, err := strconv.Atoi(getEnv("MAX_FILE_SIZE_MB", "5"))
	if err != nil {
		maxFileSize = 5
	}

	AppConfig = &Config{
		AppName:        getEnv("APP_NAME", "IT Helpdesk Ticketing API"),
		AppPort:        getEnv("APP_PORT", "8086"),
		AppEnv:         getEnv("APP_ENV", "development"),
		AppURL:         getEnv("APP_URL", "http://localhost:8086"),
		DBHost:         getEnv("DB_HOST", "127.0.0.1"),
		DBPort:         getEnv("DB_PORT", "5432"),
		DBUser:         getEnv("DB_USER", "postgres"),
		DBPassword:     getEnv("DB_PASSWORD", ""),
		DBName:         getEnv("DB_NAME", "ticketing_db"),
		DBSSLMode:      getEnv("DB_SSLMODE", "disable"),
		DBTimeZone:     getEnv("DB_TIMEZONE", "Asia/Jakarta"),
		JWTSecret:      getEnv("JWT_SECRET", "default_secret_key_123"),
		JWTExpireHours: jwtExpire,
		UploadDir:      getEnv("UPLOAD_DIR", "./storage/uploads/tickets"),
		MaxFileSizeMB:  maxFileSize,
	}

	return AppConfig
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

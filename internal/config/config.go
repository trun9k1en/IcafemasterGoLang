package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	Server  ServerConfig
	MongoDB MongoDBConfig
	Upload  UploadConfig
	JWT     JWTConfig
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	SecretKey            string
	AccessTokenDuration  int64 // in minutes
	RefreshTokenDuration int64 // in hours
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port string
	Host string
}

// MongoDBConfig holds MongoDB configuration
type MongoDBConfig struct {
	URI      string
	Database string
}

// UploadConfig holds file upload configuration
type UploadConfig struct {
	Path         string
	MaxFileSize  int64
	AllowedTypes []string
	BaseURL      string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	// Load .env file if exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	maxFileSize, _ := strconv.ParseInt(getEnv("MAX_FILE_SIZE", "52428800"), 10, 64)                  // 50MB default
	accessTokenDuration, _ := strconv.ParseInt(getEnv("JWT_ACCESS_TOKEN_DURATION", "15"), 10, 64)    // 15 minutes
	refreshTokenDuration, _ := strconv.ParseInt(getEnv("JWT_REFRESH_TOKEN_DURATION", "168"), 10, 64) // 7 days

	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
		},
		MongoDB: MongoDBConfig{
			URI:      getEnv("MONGODB_URI", "mongodb://localhost:27017"),
			Database: getEnv("MONGODB_DATABASE", "icafe_registration"),
		},
		Upload: UploadConfig{
			Path: getEnv("UPLOAD_PATH", "uploads"), // ✅ KHÔNG ./

			MaxFileSize: maxFileSize,
			AllowedTypes: []string{
				"image/jpeg", "image/png", "image/gif",
				"video/mp4", "video/mpeg", "video/quicktime", "video/webm",
				"application/pdf", "application/zip",
				"application/msword",
				"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
				"application/vnd.android.package-archive", // Cho file .apk
				"application/x-msdownload",                // Cho file .exe
				"application/octet-stream",
			},
			BaseURL: getEnv("BASE_URL", "http://localhost:8080"),
		},
		JWT: JWTConfig{
			SecretKey:            getEnv("JWT_SECRET_KEY", "your-super-secret-key-change-in-production"),
			AccessTokenDuration:  accessTokenDuration,
			RefreshTokenDuration: refreshTokenDuration,
		},
	}
}

// getEnv gets environment variable with default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

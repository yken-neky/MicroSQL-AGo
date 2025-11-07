package config

import (
	"os"
)

// Config holds basic configuration for the application
type Config struct {
	ServerPort string
	DBPath     string
	JWTSecret  string
	EncKey     string
	LogLevel   string
	// MSSQL settings
	MssqlHost string
	MssqlPort string
	MssqlUser string
	MssqlPass string
	MssqlDB   string
}

// LoadConfig loads configuration from environment variables with sensible defaults
func LoadConfig() Config {
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8000"
	}
	db := os.Getenv("DB_PATH")
	if db == "" {
		db = "./db.sqlite3"
	}
	jwt := os.Getenv("JWT_SECRET")
	if jwt == "" {
		jwt = "change-me-in-production"
	}
	enc := os.Getenv("ENCRYPTION_KEY")
	if enc == "" {
		enc = "01234567890123456789012345678901" // 32 bytes default (dev only)
	}
	loglevel := os.Getenv("LOG_LEVEL")
	if loglevel == "" {
		loglevel = "info"
	}

	return Config{
		ServerPort: port,
		DBPath:     db,
		JWTSecret:  jwt,
		EncKey:     enc,
		LogLevel:   loglevel,
		MssqlHost:  os.Getenv("MSSQL_HOST"),
		MssqlPort:  os.Getenv("MSSQL_PORT"),
		MssqlUser:  os.Getenv("MSSQL_USER"),
		MssqlPass:  os.Getenv("MSSQL_PASSWORD"),
		MssqlDB:    os.Getenv("MSSQL_DATABASE"),
	}
}

// NewGormDB simplified helper - placed here for quick access
// The actual DB creation is implemented in database.go to keep single responsibility

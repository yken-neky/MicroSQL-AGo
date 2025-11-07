package config

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// NewGormDB opens a GORM DB connection using the provided Config.
// If MySQL environment variables are present the function will open a
// MySQL connection, otherwise it falls back to SQLite using DBPath.
func NewGormDB(cfg Config) (*gorm.DB, error) {
	// Prefer MySQL when configured
	if cfg.MysqlHost != "" && cfg.MysqlUser != "" {
		port := cfg.MysqlPort
		if port == "" {
			port = "3306"
		}
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.MysqlUser, cfg.MysqlPass, cfg.MysqlHost, port, cfg.MysqlDB,
		)
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			return nil, fmt.Errorf("failed to open mysql db: %w", err)
		}
		return db, nil
	}

	// Fallback to sqlite (local file)
	db, err := gorm.Open(sqlite.Open(cfg.DBPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

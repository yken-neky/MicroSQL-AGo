package config

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// NewGormDB opens a GORM DB connection using the provided Config
func NewGormDB(cfg Config) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(cfg.DBPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

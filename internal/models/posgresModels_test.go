package models

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestAutoMigrate(t *testing.T) {
	// Create an in-memory SQLite database for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	// Run the AutoMigrate function
	err = AutoMigrate(db)
	if err != nil {
		t.Fatalf("AutoMigrate failed: %v", err)
	}

	// Check if the table was created
	if !db.Migrator().HasTable(&Control{}) {
		t.Errorf("expected table 'controls' to be created")
	}
}

package connections

import (
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestSetupDatabase_Success(t *testing.T) {
	// Mock the DSN for testing
	dsn := `host=localhost 
			user=postgres 
			password=POSTGRE*SQL 
			dbname=MicroSQL_AGo 
			port=5432 
			sslmode=disable`

	// Open a connection to the database
	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	// Check if the connection is not nil
	if DB == nil {
		t.Fatal("Expected a valid database connection, got nil")
	}

	// Check if the connection pool is configured correctly
	sqlDB, err := DB.DB()
	if err != nil {
		t.Fatalf("Failed to get sqlDB from gorm DB: %v", err)
	}

	maxIdleConns := sqlDB.Stats().Idle
	if maxIdleConns != 10 {
		t.Errorf("Expected max idle connections to be 10, got %d", maxIdleConns)
	}

	maxOpenConns := sqlDB.Stats().OpenConnections
	if maxOpenConns != 100 {
		t.Errorf("Expected max open connections to be 100, got %d", maxOpenConns)
	}
}

func TestSetupDatabase_Failure(t *testing.T) {
	// Mock an invalid DSN for testing failure
	dsn := `host=invalidhost 
			user=invaliduser 
			password=invalidpassword 
			dbname=invalidDB 
			port=5432 
			sslmode=disable`

	// Attempt to open a connection to the database
	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err == nil {
		t.Fatal("Expected connection to fail, but it succeeded")
	}

	// Check if the connection is nil
	if DB != nil {
		t.Fatal("Expected a nil database connection, got a valid connection")
	}
}

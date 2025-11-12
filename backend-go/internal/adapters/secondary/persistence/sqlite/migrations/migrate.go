package migrations

import (
	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/entities"
	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/infrastructure/repositories"
	"gorm.io/gorm"
)

// Migrate runs GORM AutoMigrate for domain entities and repository DB structs
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&entities.User{},
		&entities.ActiveConnection{},
		&entities.ControlsInformation{},
		&entities.Role{},
		&entities.Permission{},
		&entities.UserRole{},
		&entities.ConnectionLog{},
		&entities.Query{},
		&entities.ExecutionStats{},
		&repositories.QueryResultDB{},
	)
}

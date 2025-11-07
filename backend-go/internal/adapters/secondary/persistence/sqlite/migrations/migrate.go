package migrations

import (
	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/entities"
	"gorm.io/gorm"
)

// Migrate runs GORM AutoMigrate for domain entities
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&entities.User{},
		&entities.ActiveConnection{},
		&entities.ControlsInformation{},
	)
}

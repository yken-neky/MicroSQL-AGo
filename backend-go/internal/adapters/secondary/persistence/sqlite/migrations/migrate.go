package migrations

import (
	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/entities"
	// (no repository DB structs required here after removing query-related models)
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
		// Query-related models were removed (no user-facing SQL execution/persistence)
		&entities.AuditRun{},
		&entities.AuditScriptResult{},
		&entities.AdminActionLog{},
		&entities.Session{},
	)
}

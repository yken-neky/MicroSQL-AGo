package repositories

import (
	"fmt"

	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/entities"
	"gorm.io/gorm"
)

// SeedDefaultRolesAndPermissions inserts a small set of initial roles and permissions
// if they do not already exist. This is safe to call multiple times.
func SeedDefaultRolesAndPermissions(db *gorm.DB) error {
	// Ensure tables exist
	if err := db.AutoMigrate(&entities.Role{}, &entities.Permission{}, &entities.UserRole{}); err != nil {
		return fmt.Errorf("migrate roles/permissions: %w", err)
	}

	// default roles
	roles := []entities.Role{
		{Name: "admin", Description: "Full system administrator"},
		{Name: "auditor", Description: "Can run and view audits"},
		{Name: "user", Description: "Regular user with limited privileges"},
	}
	for _, r := range roles {
		var existing entities.Role
		if err := db.Where("name = ?", r.Name).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&r).Error; err != nil {
					return fmt.Errorf("create role %s: %w", r.Name, err)
				}
			} else {
				return err
			}
		}
	}

	perms := []entities.Permission{
		{Name: "users:create", Resource: "users", Action: "create", Description: "Create users"},
		{Name: "users:read", Resource: "users", Action: "read", Description: "Read users"},
		{Name: "users:update", Resource: "users", Action: "update", Description: "Update users"},
		{Name: "users:delete", Resource: "users", Action: "delete", Description: "Delete users"},
		{Name: "connections:manage", Resource: "connections", Action: "manage", Description: "Manage DB connections"},
		{Name: "audits:execute", Resource: "audits", Action: "execute", Description: "Execute audit scripts"},
		{Name: "audits:view", Resource: "audits", Action: "view", Description: "View audit results"},
		{Name: "roles:manage", Resource: "roles", Action: "manage", Description: "Manage roles and assignments"},
		{Name: "permissions:manage", Resource: "permissions", Action: "manage", Description: "Manage permissions"},
	}

	for _, p := range perms {
		var existing entities.Permission
		if err := db.Where("name = ?", p.Name).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&p).Error; err != nil {
					return fmt.Errorf("create permission %s: %w", p.Name, err)
				}
			} else {
				return err
			}
		}
	}

	// Attach sensible defaults: admin gets all permissions; auditor gets audits:view/execute; user has none by default
	var adminRole entities.Role
	if err := db.Preload("Permissions").Where("name = ?", "admin").First(&adminRole).Error; err != nil {
		return err
	}

	var allPerms []entities.Permission
	if err := db.Find(&allPerms).Error; err != nil {
		return err
	}

	if err := db.Model(&adminRole).Association("Permissions").Replace(&allPerms); err != nil {
		return fmt.Errorf("assign perms to admin: %w", err)
	}

	var auditorRole entities.Role
	if err := db.Preload("Permissions").Where("name = ?", "auditor").First(&auditorRole).Error; err != nil {
		return err
	}

	var audPerms []entities.Permission
	if err := db.Where("resource = ?", "audits").Find(&audPerms).Error; err != nil {
		return err
	}
	if err := db.Model(&auditorRole).Association("Permissions").Replace(&audPerms); err != nil {
		return fmt.Errorf("assign perms to auditor: %w", err)
	}

	return nil
}

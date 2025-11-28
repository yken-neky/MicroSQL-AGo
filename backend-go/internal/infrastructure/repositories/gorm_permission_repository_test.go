package repositories

import (
	"testing"

	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/entities"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestGormPermissionRepository_CRUD_And_AssignToRole(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite: %v", err)
	}

	if err := db.AutoMigrate(&entities.Role{}, &entities.Permission{}); err != nil {
		t.Fatalf("migrate failed: %v", err)
	}

	repo := NewGormPermissionRepository(db)

	perm := &entities.Permission{Name: "users:read", Resource: "users", Action: "read", Description: "read users"}
	if err := repo.Create(perm); err != nil {
		t.Fatalf("create permission: %v", err)
	}
	if perm.ID == 0 {
		t.Fatalf("expected permission ID to be set")
	}

	// Create role and assign permission
	role := &entities.Role{Name: "auditor", Description: "auditor role"}
	if err := db.Create(role).Error; err != nil {
		t.Fatalf("create role: %v", err)
	}

	if err := repo.AssignToRole(role.ID, perm.ID); err != nil {
		t.Fatalf("assign permission to role: %v", err)
	}

	list, err := repo.GetRolePermissions(role.ID)
	if err != nil {
		t.Fatalf("get role permissions: %v", err)
	}
	if len(list) != 1 || list[0].Name != "users:read" {
		t.Fatalf("expected permission assigned to role")
	}

	ok, err := repo.HasPermission(role.ID, "users:read")
	if err != nil {
		t.Fatalf("has permission check: %v", err)
	}
	if !ok {
		t.Fatalf("expected HasPermission true")
	}

	// Revoke
	if err := repo.RevokeFromRole(role.ID, perm.ID); err != nil {
		t.Fatalf("revoke permission: %v", err)
	}

	perms2, err := repo.GetRolePermissions(role.ID)
	if err != nil {
		t.Fatalf("get role perms after revoke: %v", err)
	}
	if len(perms2) != 0 {
		t.Fatalf("expected no perms after revoke")
	}
}

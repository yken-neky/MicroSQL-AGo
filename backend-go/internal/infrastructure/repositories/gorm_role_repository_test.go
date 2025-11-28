package repositories

import (
	"testing"

	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/entities"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestGormRoleRepository_CRUD_And_UserAssignment(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite: %v", err)
	}

	if err := db.AutoMigrate(&entities.User{}, &entities.Role{}, &entities.UserRole{}, &entities.Permission{}); err != nil {
		t.Fatalf("migrate failed: %v", err)
	}

	repo := NewGormRoleRepository(db)

	role := &entities.Role{Name: "tester", Description: "for tests"}
	if err := repo.Create(role); err != nil {
		t.Fatalf("create role: %v", err)
	}
	if role.ID == 0 {
		t.Fatalf("expected role ID to be set")
	}

	// GetByName
	got, err := repo.GetByName("tester")
	if err != nil {
		t.Fatalf("get by name: %v", err)
	}
	if got == nil || got.Name != "tester" {
		t.Fatalf("unexpected get by name result")
	}

	// Update
	role.Description = "updated"
	if err := repo.Update(role); err != nil {
		t.Fatalf("update role: %v", err)
	}

	// Assign to user
	user := &entities.User{Username: "u1", Email: "u1@example.com", Password: "x"}
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}

	if err := repo.AssignToUser(user.ID, role.ID); err != nil {
		t.Fatalf("assign role to user: %v", err)
	}

	roles, err := repo.GetUserRoles(user.ID)
	if err != nil {
		t.Fatalf("get user roles: %v", err)
	}
	if len(roles) != 1 || roles[0].Name != "tester" {
		t.Fatalf("expected assigned role present")
	}

	// Revoke
	if err := repo.RevokeFromUser(user.ID, role.ID); err != nil {
		t.Fatalf("revoke role: %v", err)
	}

	roles2, err := repo.GetUserRoles(user.ID)
	if err != nil {
		t.Fatalf("get user roles after revoke: %v", err)
	}
	if len(roles2) != 0 {
		t.Fatalf("expected no roles after revoke")
	}

	// Delete
	if err := repo.Delete(role.ID); err != nil {
		t.Fatalf("delete role: %v", err)
	}
}

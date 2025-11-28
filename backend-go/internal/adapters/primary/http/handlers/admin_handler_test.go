package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/entities"
	repo "github.com/yken-neky/MicroSQL-AGo/backend-go/internal/infrastructure/repositories"
)

func setupDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	return db
}

func TestAdminHandler_GetUsersMetrics(t *testing.T) {
	db := setupDB(t)
	if err := db.AutoMigrate(&entities.User{}); err != nil {
		t.Fatalf("migrate users: %v", err)
	}

	// seed users
	db.Create(&entities.User{Username: "a", Email: "a@e", Password: "x", Role: "admin", IsActive: true})
	db.Create(&entities.User{Username: "b", Email: "b@e", Password: "x", Role: "user", IsActive: false})
	db.Create(&entities.User{Username: "c", Email: "c@e", Password: "x", Role: "user", IsActive: true})

	roleRepo := repo.NewGormRoleRepository(db)
	permRepo := repo.NewGormPermissionRepository(db)
	auditRepo := repo.NewGormAdminAuditRepository(db)
	h := NewAdminHandler(db, zap.NewNop(), nil, roleRepo, permRepo, auditRepo)

	r := gin.New()
	// test routes simulate authenticated admin user in context
	r.Use(func(c *gin.Context) { c.Set("userID", uint(99)); c.Set("username", "test-admin"); c.Next() })
	r.GET("/metrics/users", h.GetUsersMetrics)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/metrics/users", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d body=%s", w.Code, w.Body.String())
	}
}

func TestAdminHandler_GetConnectionsMetrics(t *testing.T) {
	db := setupDB(t)
	if err := db.AutoMigrate(&entities.ActiveConnection{}, &entities.ConnectionLog{}); err != nil {
		t.Fatalf("migrate conns: %v", err)
	}

	db.Create(&entities.ActiveConnection{Driver: "mssql", Server: "s1", DBUser: "u1", IsConnected: true})
	db.Create(&entities.ActiveConnection{Driver: "mssql", Server: "s2", DBUser: "u2", IsConnected: false})
	db.Create(&entities.ConnectionLog{Driver: "mssql", Server: "s1", DBUser: "u1", Status: "ok", UserID: 1})

	roleRepo := repo.NewGormRoleRepository(db)
	permRepo := repo.NewGormPermissionRepository(db)
	auditRepo := repo.NewGormAdminAuditRepository(db)
	h := NewAdminHandler(db, zap.NewNop(), nil, roleRepo, permRepo, auditRepo)
	r := gin.New()
	r.Use(func(c *gin.Context) { c.Set("userID", uint(99)); c.Set("username", "test-admin"); c.Next() })
	r.GET("/metrics/connections", h.GetConnectionsMetrics)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/metrics/connections", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d body=%s", w.Code, w.Body.String())
	}
}

func TestAdminHandler_GetAuditsMetrics(t *testing.T) {
	db := setupDB(t)
	if err := db.AutoMigrate(&entities.AuditRun{}, &entities.AuditScriptResult{}); err != nil {
		t.Fatalf("migrate audits: %v", err)
	}

	now := time.Now()
	finished := now.Add(2 * time.Second)
	db.Create(&entities.AuditRun{UserID: 1, Mode: "partial", Status: "running", StartedAt: now})
	db.Create(&entities.AuditRun{UserID: 1, Mode: "partial", Status: "completed", StartedAt: now.Add(-10 * time.Second), FinishedAt: &finished, Total: 1, Passed: 1})

	roleRepo := repo.NewGormRoleRepository(db)
	permRepo := repo.NewGormPermissionRepository(db)
	auditRepo := repo.NewGormAdminAuditRepository(db)
	h := NewAdminHandler(db, zap.NewNop(), nil, roleRepo, permRepo, auditRepo)
	r := gin.New()
	r.Use(func(c *gin.Context) { c.Set("userID", uint(99)); c.Set("username", "test-admin"); c.Next() })
	r.GET("/metrics/audits", h.GetAuditsMetrics)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/metrics/audits", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d body=%s", w.Code, w.Body.String())
	}
}

func TestAdminHandler_GetRolesMetricsAndSystemMetrics(t *testing.T) {
	db := setupDB(t)
	if err := db.AutoMigrate(&entities.Role{}, &entities.Permission{}, &entities.UserRole{}); err != nil {
		t.Fatalf("migrate roles: %v", err)
	}

	// create a role and permission and associate
	role := entities.Role{Name: "admin"}
	perm := entities.Permission{Name: "users:read", Resource: "users", Action: "read"}
	db.Create(&role)
	db.Create(&perm)
	db.Model(&role).Association("Permissions").Append(&perm)

	roleRepo := repo.NewGormRoleRepository(db)
	permRepo := repo.NewGormPermissionRepository(db)
	auditRepo := repo.NewGormAdminAuditRepository(db)
	h := NewAdminHandler(db, zap.NewNop(), nil, roleRepo, permRepo, auditRepo)
	r := gin.New()
	r.Use(func(c *gin.Context) { c.Set("userID", uint(99)); c.Set("username", "test-admin"); c.Next() })
	r.GET("/metrics/roles", h.GetRolesMetrics)
	r.GET("/metrics/system", h.GetSystemMetrics)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/metrics/roles", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d body=%s", w.Code, w.Body.String())
	}

	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodGet, "/metrics/system", nil)
	r.ServeHTTP(w2, req2)
	if w2.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d body=%s", w2.Code, w2.Body.String())
	}
}

func TestAdminHandler_RolesAndPermissionsCRUD(t *testing.T) {
	db := setupDB(t)
	if err := db.AutoMigrate(&entities.User{}, &entities.Role{}, &entities.Permission{}, &entities.UserRole{}, &entities.AdminActionLog{}); err != nil {
		t.Fatalf("migrate roles/permissions: %v", err)
	}

	roleRepo := repo.NewGormRoleRepository(db)
	permRepo := repo.NewGormPermissionRepository(db)
	auditRepo := repo.NewGormAdminAuditRepository(db)
	h := NewAdminHandler(db, zap.NewNop(), nil, roleRepo, permRepo, auditRepo)

	r := gin.New()
	r.Use(func(c *gin.Context) { c.Set("userID", uint(99)); c.Set("username", "test-admin"); c.Next() })
	r.POST("/roles", h.CreateRole)
	r.GET("/roles", h.ListRoles)
	r.PUT("/roles/:id", h.UpdateRole)
	r.DELETE("/roles/:id", h.DeleteRole)

	// create role
	createReq := httptest.NewRequest(http.MethodPost, "/roles", strings.NewReader(`{"name":"dev","description":"dev role"}`))
	createReq.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, createReq)
	if w.Code != http.StatusCreated {
		t.Fatalf("create role failed: %d body=%s", w.Code, w.Body.String())
	}

	// list roles
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, httptest.NewRequest(http.MethodGet, "/roles", nil))
	if w2.Code != http.StatusOK {
		t.Fatalf("list roles failed: %d body=%s", w2.Code, w2.Body.String())
	}

	// parse created role id
	var out struct {
		Roles []entities.Role `json:"roles"`
	}
	if err := json.Unmarshal(w2.Body.Bytes(), &out); err != nil || len(out.Roles) == 0 {
		t.Fatalf("unexpected list output: %v err=%v", w2.Body.String(), err)
	}
	rid := out.Roles[0].ID

	// update
	updReq := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/roles/%d", rid), strings.NewReader(`{"description":"updated"}`))
	updReq.Header.Set("Content-Type", "application/json")
	w3 := httptest.NewRecorder()
	r.ServeHTTP(w3, updReq)
	if w3.Code != http.StatusOK {
		t.Fatalf("update role failed: %d body=%s", w3.Code, w3.Body.String())
	}

	// delete
	w4 := httptest.NewRecorder()
	r.ServeHTTP(w4, httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/roles/%d", rid), nil))
	if w4.Code != http.StatusNoContent {
		t.Fatalf("delete role failed: %d body=%s", w4.Code, w4.Body.String())
	}

	// Permissions CRUD + assign
	r2 := gin.New()
	r2.POST("/permissions", h.CreatePermission)
	r2.GET("/permissions", h.ListPermissions)
	r2.PUT("/permissions/:id", h.UpdatePermission)
	r2.DELETE("/permissions/:id", h.DeletePermission)
	r2.POST("/roles/:id/permissions", h.AssignPermissionToRole)

	// create permission
	preq := httptest.NewRequest(http.MethodPost, "/permissions", strings.NewReader(`{"name":"x:read","resource":"x","action":"read","description":"test"}`))
	preq.Header.Set("Content-Type", "application/json")
	pw := httptest.NewRecorder()
	r2.ServeHTTP(pw, preq)
	if pw.Code != http.StatusCreated {
		t.Fatalf("create perm failed: %d body=%s", pw.Code, pw.Body.String())
	}

	// list perms
	pw2 := httptest.NewRecorder()
	r2.ServeHTTP(pw2, httptest.NewRequest(http.MethodGet, "/permissions", nil))
	if pw2.Code != http.StatusOK {
		t.Fatalf("list perms failed: %d body=%s", pw2.Code, pw2.Body.String())
	}
	var outp struct {
		Permissions []entities.Permission `json:"permissions"`
	}
	if err := json.Unmarshal(pw2.Body.Bytes(), &outp); err != nil || len(outp.Permissions) == 0 {
		t.Fatalf("unexpected perms list: %v err=%v", pw2.Body.String(), err)
	}
	pid := outp.Permissions[0].ID

	// assign permission to a role
	// create a role first
	rr := repo.NewGormRoleRepository(db)
	rrole := entities.Role{Name: "assignable"}
	rr.Create(&rrole)

	assignReq := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/roles/%d/permissions", rrole.ID), strings.NewReader(fmt.Sprintf(`{"permission_id": %d}`, pid)))
	assignReq.Header.Set("Content-Type", "application/json")
	aw := httptest.NewRecorder()
	r2.ServeHTTP(aw, assignReq)
	if aw.Code != http.StatusOK {
		t.Fatalf("assign perm failed: %d body=%s", aw.Code, aw.Body.String())
	}

	// verify RBAC audit logs were recorded
	logs, err := auditRepo.List(nil, nil, nil, 100, 0)
	if err != nil {
		t.Fatalf("failed reading audit logs: %v", err)
	}
	if len(logs) == 0 {
		t.Fatalf("expected audit logs to contain entries for RBAC actions, got 0")
	}
}

func TestAdminHandler_ListRBACAuditLogs(t *testing.T) {
	db := setupDB(t)
	if err := db.AutoMigrate(&entities.AdminActionLog{}); err != nil {
		t.Fatalf("migrate audit logs: %v", err)
	}

	auditRepo := repo.NewGormAdminAuditRepository(db)
	// seed some audit entries
	auditRepo.Create(&entities.AdminActionLog{ActorID: 1, ActorName: "admin1", Action: "role.create", TargetType: "role", TargetName: "admin", Details: "created"})
	auditRepo.Create(&entities.AdminActionLog{ActorID: 2, ActorName: "auditor", Action: "permission.create", TargetType: "permission", TargetName: "audits:read", Details: "added perm"})

	roleRepo := repo.NewGormRoleRepository(db)
	permRepo := repo.NewGormPermissionRepository(db)
	h := NewAdminHandler(db, zap.NewNop(), nil, roleRepo, permRepo, auditRepo)

	r := gin.New()
	r.Use(func(c *gin.Context) { c.Set("userID", uint(99)); c.Set("username", "test-admin"); c.Next() })
	r.GET("/admin/audit/rbac", h.ListRBACAuditLogs)

	// list all
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/admin/audit/rbac", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d body=%s", w.Code, w.Body.String())
	}

	var out struct {
		Logs []entities.AdminActionLog `json:"logs"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &out); err != nil {
		t.Fatalf("failed to parse response: %v body=%s", err, w.Body.String())
	}
	if len(out.Logs) != 2 {
		t.Fatalf("expected 2 logs, got %d", len(out.Logs))
	}

	// filter by actor id
	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodGet, "/admin/audit/rbac?actor_id=1", nil)
	r.ServeHTTP(w2, req2)
	if w2.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d body=%s", w2.Code, w2.Body.String())
	}
	var out2 struct {
		Logs []entities.AdminActionLog `json:"logs"`
	}
	if err := json.Unmarshal(w2.Body.Bytes(), &out2); err != nil {
		t.Fatalf("failed to parse response: %v body=%s", err, w2.Body.String())
	}
	if len(out2.Logs) != 1 || out2.Logs[0].ActorID != 1 {
		t.Fatalf("expected single log for actor_id=1, got %+v", out2.Logs)
	}
}

// ...existing code...
func TestAdminHandler_ListUsersWithRoles(t *testing.T) {
    db := setupDB(t)
    if err := db.AutoMigrate(&entities.User{}, &entities.Role{}, &entities.UserRole{}); err != nil {
        t.Fatalf("migrate users/roles: %v", err)
    }

    // create roles
    roleAdmin := entities.Role{Name: "admin"}
    roleManager := entities.Role{Name: "manager"}
    db.Create(&roleAdmin)
    db.Create(&roleManager)

    // create users
    u1 := entities.User{Username: "alice", Email: "alice@e", Password: "x", IsActive: true}
    u2 := entities.User{Username: "bob", Email: "bob@e", Password: "x", IsActive: true}
    db.Create(&u1)
    db.Create(&u2)

    // assign roles via user_roles table
    db.Create(&entities.UserRole{UserID: u1.ID, RoleID: roleAdmin.ID})
    db.Create(&entities.UserRole{UserID: u2.ID, RoleID: roleManager.ID})
    db.Create(&entities.UserRole{UserID: u2.ID, RoleID: roleAdmin.ID}) // bob has two roles

    roleRepo := repo.NewGormRoleRepository(db)
    permRepo := repo.NewGormPermissionRepository(db)
    auditRepo := repo.NewGormAdminAuditRepository(db)
    h := NewAdminHandler(db, zap.NewNop(), nil, roleRepo, permRepo, auditRepo)

    r := gin.New()
    r.Use(func(c *gin.Context) { c.Set("userID", uint(99)); c.Set("username", "test-admin"); c.Next() })
    r.GET("/admin/users", h.ListUsersWithRoles)

    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodGet, "/admin/users", nil)
    r.ServeHTTP(w, req)

    if w.Code != http.StatusOK {
        t.Fatalf("expected 200 got %d body=%s", w.Code, w.Body.String())
    }

    var out struct {
        Users []struct {
            ID       uint     `json:"id"`
            Username string   `json:"username"`
            Email    string   `json:"email"`
            IsActive bool     `json:"is_active"`
            Roles    []string `json:"roles"`
        } `json:"users"`
    }
    if err := json.Unmarshal(w.Body.Bytes(), &out); err != nil {
        t.Fatalf("failed to parse response: %v body=%s", err, w.Body.String())
    }

    if len(out.Users) != 2 {
        t.Fatalf("expected 2 users, got %d", len(out.Users))
    }

    // verify roles returned for users
    for _, u := range out.Users {
        switch u.Username {
        case "alice":
            if len(u.Roles) != 1 || u.Roles[0] != "admin" {
                t.Fatalf("alice roles mismatch: %+v", u.Roles)
            }
        case "bob":
            if len(u.Roles) != 2 {
                t.Fatalf("bob should have 2 roles, got: %+v", u.Roles)
            }
        default:
            t.Fatalf("unexpected user %s", u.Username)
        }
    }
}

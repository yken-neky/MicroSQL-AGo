package handlers

import (
    "net/http"
    "net/http/httptest"
    "testing"
    "time"

    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"

    "github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/entities"
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

    h := NewAdminHandler(db, zap.NewNop(), nil)

    r := gin.New()
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

    h := NewAdminHandler(db, zap.NewNop(), nil)
    r := gin.New()
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

    h := NewAdminHandler(db, zap.NewNop(), nil)
    r := gin.New()
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

    h := NewAdminHandler(db, zap.NewNop(), nil)
    r := gin.New()
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

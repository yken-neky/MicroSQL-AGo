package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/entities"
	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/ports/repositories"
	"gorm.io/gorm"
)

// AdminHandler for admin-only endpoints
type AdminHandler struct {
	DB          *gorm.DB
	Logger      *zap.Logger
	SessionRepo repositories.SessionRepository
	RoleRepo    repositories.RoleRepository
	PermRepo    repositories.PermissionRepository
}

func NewAdminHandler(db *gorm.DB, logger *zap.Logger, sr repositories.SessionRepository, rr repositories.RoleRepository, pr repositories.PermissionRepository) *AdminHandler {
	return &AdminHandler{DB: db, Logger: logger, SessionRepo: sr, RoleRepo: rr, PermRepo: pr}
}

// ListActiveSessions returns all active sessions with user info (admin only)
func (h *AdminHandler) ListActiveSessions(c *gin.Context) {
	if h.SessionRepo == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "session repository not configured"})
		return
	}

	sessions, err := h.SessionRepo.ListActiveSessions()
	if err != nil {
		h.Logger.Error("failed listing active sessions", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list sessions"})
		return
	}

	resp := make([]map[string]interface{}, 0, len(sessions))
	for _, s := range sessions {
		// fetch user info
		var u entities.User
		if err := h.DB.First(&u, s.UserID).Error; err != nil {
			// if user not found, still include session, with null user
			resp = append(resp, map[string]interface{}{
				"session_id": s.ID,
				"user_id":    s.UserID,
				"username":   nil,
				"email":      nil,
				"token":      s.Token,
				"expires_at": s.ExpiresAt,
				"created_at": s.CreatedAt,
			})
			continue
		}
		resp = append(resp, map[string]interface{}{
			"session_id": s.ID,
			"user_id":    s.UserID,
			"username":   u.Username,
			"email":      u.Email,
			"token":      s.Token,
			"expires_at": s.ExpiresAt,
			"created_at": s.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"sessions": resp})
}

// GetUsersMetrics returns basic user metrics (counts and distribution by role)
func (h *AdminHandler) GetUsersMetrics(c *gin.Context) {
	var total int64
	if err := h.DB.Model(&entities.User{}).Count(&total).Error; err != nil {
		h.Logger.Error("count users failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to collect user metrics"})
		return
	}

	var active int64
	if err := h.DB.Model(&entities.User{}).Where("is_active = ?", true).Count(&active).Error; err != nil {
		h.Logger.Error("count active users failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to collect user metrics"})
		return
	}

	// roles distribution
	type RoleCount struct {
		Role  string
		Count int64
	}
	var roleCounts []RoleCount
	if err := h.DB.Model(&entities.User{}).Select("role, COUNT(*) as count").Group("role").Scan(&roleCounts).Error; err != nil {
		h.Logger.Error("role distribution failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to collect role distribution"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"total_users": total, "active_users": active, "roles_distribution": roleCounts})
}

// GetConnectionsMetrics returns stats about connections and connection_logs
func (h *AdminHandler) GetConnectionsMetrics(c *gin.Context) {
	var totalActive int64
	if err := h.DB.Model(&entities.ActiveConnection{}).Count(&totalActive).Error; err != nil {
		h.Logger.Error("count active connections failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to collect connection metrics"})
		return
	}

	var connectedNow int64
	if err := h.DB.Model(&entities.ActiveConnection{}).Where("is_connected = ?", true).Count(&connectedNow).Error; err != nil {
		h.Logger.Error("count connected now failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to collect connection metrics"})
		return
	}

	var logsTotal int64
	if err := h.DB.Model(&entities.ConnectionLog{}).Count(&logsTotal).Error; err != nil {
		h.Logger.Error("count connection logs failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to collect connection logs metrics"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"total_active_connections": totalActive, "currently_connected": connectedNow, "total_connection_logs": logsTotal})
}

// GetAuditsMetrics returns stats about audit runs and script results
func (h *AdminHandler) GetAuditsMetrics(c *gin.Context) {
	var totalRuns int64
	if err := h.DB.Model(&entities.AuditRun{}).Count(&totalRuns).Error; err != nil {
		h.Logger.Error("count audit runs failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to collect audits metrics"})
		return
	}

	type StatusCount struct {
		Status string
		Count  int64
	}

	var statusCounts []StatusCount
	if err := h.DB.Model(&entities.AuditRun{}).Select("status, COUNT(*) as count").Group("status").Scan(&statusCounts).Error; err != nil {
		h.Logger.Error("audit status distribution failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to collect audit status distribution"})
		return
	}

	// average duration of completed runs (in seconds)
	var runs []entities.AuditRun
	if err := h.DB.Where("finished_at IS NOT NULL").Find(&runs).Error; err != nil {
		h.Logger.Error("reading audit durations failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to compute audits duration"})
		return
	}

	var totalDuration float64
	var completedCount int64
	for _, r := range runs {
		if r.FinishedAt != nil {
			totalDuration += r.FinishedAt.Sub(r.StartedAt).Seconds()
			completedCount++
		}
	}

	avg := 0.0
	if completedCount > 0 {
		avg = totalDuration / float64(completedCount)
	}

	c.JSON(http.StatusOK, gin.H{"total_runs": totalRuns, "status_distribution": statusCounts, "average_duration_seconds": avg})
}

// GetRolesMetrics returns counts for roles and permissions
func (h *AdminHandler) GetRolesMetrics(c *gin.Context) {
	var totalRoles int64
	if err := h.DB.Model(&entities.Role{}).Count(&totalRoles).Error; err != nil {
		h.Logger.Error("count roles failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to collect roles metrics"})
		return
	}

	var totalPerms int64
	if err := h.DB.Model(&entities.Permission{}).Count(&totalPerms).Error; err != nil {
		h.Logger.Error("count permissions failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to collect permissions metrics"})
		return
	}

	// permissions per role
	type RolePermCount struct {
		Role  string
		Count int64
	}
	var rpc []RolePermCount
	// join to role_permissions
	if err := h.DB.Table("roles").Select("roles.name as role, COUNT(role_permissions.permission_id) as count").Joins("LEFT JOIN role_permissions ON role_permissions.role_id = roles.id").Group("roles.id").Scan(&rpc).Error; err != nil {
		h.Logger.Error("role-permission scan failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to collect role-permission distribution"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"total_roles": totalRoles, "total_permissions": totalPerms, "permissions_per_role": rpc})
}

// GetSystemMetrics returns simple counts for important tables
func (h *AdminHandler) GetSystemMetrics(c *gin.Context) {
	tables := []string{"users", "active_connections", "connection_logs", "controls_informations", "sessions", "audit_runs", "audit_script_results", "roles", "permissions", "user_roles"}
	resp := map[string]int64{}
	for _, t := range tables {
		var cnt int64
		if err := h.DB.Table(t).Count(&cnt).Error; err != nil {
			// don't fail the whole endpoint on one table; record -1 for errors
			resp[t] = -1
			continue
		}
		resp[t] = cnt
	}
	c.JSON(http.StatusOK, gin.H{"table_counts": resp})
}

// --- Roles / Permissions CRUD and assignments ---

// ListRoles returns available roles with permissions
func (h *AdminHandler) ListRoles(c *gin.Context) {
	if h.RoleRepo == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "role repository not configured"})
		return
	}
	list, err := h.RoleRepo.List()
	if err != nil {
		h.Logger.Error("list roles failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list roles"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"roles": list})
}

// CreateRole creates a new role
func (h *AdminHandler) CreateRole(c *gin.Context) {
	if h.RoleRepo == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "role repository not configured"})
		return
	}
	var body entities.Role
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}
	if body.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name required"})
		return
	}
	if err := h.RoleRepo.Create(&body); err != nil {
		h.Logger.Error("create role failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create role"})
		return
	}
	c.JSON(http.StatusCreated, body)
}

// UpdateRole updates role metadata
func (h *AdminHandler) UpdateRole(c *gin.Context) {
	if h.RoleRepo == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "role repository not configured"})
		return
	}
	idParam := c.Param("id")
	var id uint
	if _, err := fmt.Sscanf(idParam, "%d", &id); err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role id"})
		return
	}
	existing, err := h.RoleRepo.GetByID(id)
	if err != nil {
		h.Logger.Error("get role failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load role"})
		return
	}
	if existing == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "role not found"})
		return
	}
	var body entities.Role
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}
	if body.Name != "" {
		existing.Name = body.Name
	}
	existing.Description = body.Description
	if err := h.RoleRepo.Update(existing); err != nil {
		h.Logger.Error("update role failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update role"})
		return
	}
	c.JSON(http.StatusOK, existing)
}

// DeleteRole removes a role
func (h *AdminHandler) DeleteRole(c *gin.Context) {
	if h.RoleRepo == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "role repository not configured"})
		return
	}
	idParam := c.Param("id")
	var id uint
	if _, err := fmt.Sscanf(idParam, "%d", &id); err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role id"})
		return
	}
	if err := h.RoleRepo.Delete(id); err != nil {
		h.Logger.Error("delete role failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete role"})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{})
}

// AssignRoleToUser assigns a role to a user
func (h *AdminHandler) AssignRoleToUser(c *gin.Context) {
	if h.RoleRepo == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "role repository not configured"})
		return
	}
	userIDParam := c.Param("id")
	var userID uint
	if _, err := fmt.Sscanf(userIDParam, "%d", &userID); err != nil || userID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
	var body struct {
		RoleID uint `json:"role_id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.RoleID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "role_id required"})
		return
	}
	if err := h.RoleRepo.AssignToUser(userID, body.RoleID); err != nil {
		h.Logger.Error("assign role failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to assign role"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// RevokeRoleFromUser removes a role from a user
func (h *AdminHandler) RevokeRoleFromUser(c *gin.Context) {
	if h.RoleRepo == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "role repository not configured"})
		return
	}
	userIDParam := c.Param("id")
	var userID uint
	if _, err := fmt.Sscanf(userIDParam, "%d", &userID); err != nil || userID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
	var body struct {
		RoleID uint `json:"role_id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.RoleID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "role_id required"})
		return
	}
	if err := h.RoleRepo.RevokeFromUser(userID, body.RoleID); err != nil {
		h.Logger.Error("revoke role failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to revoke role"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// ---- Permissions ----

func (h *AdminHandler) ListPermissions(c *gin.Context) {
	if h.PermRepo == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "permission repository not configured"})
		return
	}
	list, err := h.PermRepo.List()
	if err != nil {
		h.Logger.Error("list perms failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list permissions"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"permissions": list})
}

func (h *AdminHandler) CreatePermission(c *gin.Context) {
	if h.PermRepo == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "permission repository not configured"})
		return
	}
	var body entities.Permission
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}
	if body.Name == "" || body.Resource == "" || body.Action == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name, resource and action are required"})
		return
	}
	if err := h.PermRepo.Create(&body); err != nil {
		h.Logger.Error("create perm failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create permission"})
		return
	}
	c.JSON(http.StatusCreated, body)
}

func (h *AdminHandler) UpdatePermission(c *gin.Context) {
	if h.PermRepo == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "permission repository not configured"})
		return
	}
	idParam := c.Param("id")
	var id uint
	if _, err := fmt.Sscanf(idParam, "%d", &id); err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid permission id"})
		return
	}
	existing, err := h.PermRepo.GetByID(id)
	if err != nil {
		h.Logger.Error("get perm failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load permission"})
		return
	}
	if existing == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "permission not found"})
		return
	}
	var body entities.Permission
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}
	if body.Name != "" {
		existing.Name = body.Name
	}
	if body.Resource != "" {
		existing.Resource = body.Resource
	}
	if body.Action != "" {
		existing.Action = body.Action
	}
	existing.Description = body.Description
	if err := h.PermRepo.Update(existing); err != nil {
		h.Logger.Error("update perm failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update permission"})
		return
	}
	c.JSON(http.StatusOK, existing)
}

func (h *AdminHandler) DeletePermission(c *gin.Context) {
	if h.PermRepo == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "permission repository not configured"})
		return
	}
	idParam := c.Param("id")
	var id uint
	if _, err := fmt.Sscanf(idParam, "%d", &id); err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid permission id"})
		return
	}
	if err := h.PermRepo.Delete(id); err != nil {
		h.Logger.Error("delete perm failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete permission"})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{})
}

// AssignPermissionToRole attaches a permission to a role
func (h *AdminHandler) AssignPermissionToRole(c *gin.Context) {
	if h.PermRepo == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "permission repository not configured"})
		return
	}
	roleIDParam := c.Param("id")
	var roleID uint
	if _, err := fmt.Sscanf(roleIDParam, "%d", &roleID); err != nil || roleID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role id"})
		return
	}
	var body struct {
		PermissionID uint `json:"permission_id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.PermissionID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "permission_id required"})
		return
	}
	if err := h.PermRepo.AssignToRole(roleID, body.PermissionID); err != nil {
		h.Logger.Error("assign perm to role failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to assign permission to role"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *AdminHandler) RevokePermissionFromRole(c *gin.Context) {
	if h.PermRepo == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "permission repository not configured"})
		return
	}
	roleIDParam := c.Param("id")
	var roleID uint
	if _, err := fmt.Sscanf(roleIDParam, "%d", &roleID); err != nil || roleID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role id"})
		return
	}
	var body struct {
		PermissionID uint `json:"permission_id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.PermissionID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "permission_id required"})
		return
	}
	if err := h.PermRepo.RevokeFromRole(roleID, body.PermissionID); err != nil {
		h.Logger.Error("revoke perm from role failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to revoke permission from role"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

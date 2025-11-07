package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/ports/repositories"
)

// AuthorizationMiddleware maneja la verificación de permisos
type AuthorizationMiddleware struct {
	roleRepo repositories.RoleRepository
	permRepo repositories.PermissionRepository
}

func NewAuthorizationMiddleware(
	rr repositories.RoleRepository,
	pr repositories.PermissionRepository,
) *AuthorizationMiddleware {
	return &AuthorizationMiddleware{
		roleRepo: rr,
		permRepo: pr,
	}
}

// RequirePermission verifica si el usuario tiene el permiso requerido
func (m *AuthorizationMiddleware) RequirePermission(requiredPerm string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener ID del usuario del contexto (previamente establecido por el middleware de autenticación)
		userID, exists := c.Get("userID")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized: missing user context"})
			return
		}

		// Obtener roles del usuario
		roles, err := m.roleRepo.GetUserRoles(userID.(uint))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "error fetching user roles"})
			return
		}

		// Verificar si algún rol tiene el permiso requerido
		hasPermission := false
		for _, role := range roles {
			if ok, err := m.permRepo.HasPermission(role.ID, requiredPerm); err == nil && ok {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden: insufficient permissions"})
			return
		}

		c.Next()
	}
}

// RequireRole verifica si el usuario tiene el rol requerido
func (m *AuthorizationMiddleware) RequireRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized: missing user context"})
			return
		}

		roles, err := m.roleRepo.GetUserRoles(userID.(uint))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "error fetching user roles"})
			return
		}

		hasRole := false
		for _, role := range roles {
			if role.Name == requiredRole {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden: insufficient role"})
			return
		}

		c.Next()
	}
}

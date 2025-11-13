package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/adapters/secondary/security"
)

// AuthMiddleware handles JWT validation and user context
type AuthMiddleware struct {
	jwtService *security.JWTService
}

func NewAuthMiddleware(jwtService *security.JWTService) *AuthMiddleware {
	return &AuthMiddleware{jwtService: jwtService}
}

// RequireAuth verifies JWT token and sets user info in context
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "no token provided"})
			return
		}

		// Extract Bearer token
		const bearerSchema string = "Bearer "
		if len(header) < len(bearerSchema) || header[:len(bearerSchema)] != bearerSchema {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
			return
		}

		token := header[len(bearerSchema):]

		// Validate token
		claims, err := m.jwtService.ValidateToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token: " + err.Error()})
			return
		}

		// Set user info in context
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Next()
	}
}

// RequireRole ensures user has one of the required roles
func (m *AuthMiddleware) RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// First validate token (via RequireAuth)
		m.RequireAuth()(c)
		if c.IsAborted() {
			return
		}

		// Get user role from context
		userRole, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "role not found in token"})
			return
		}

		role, ok := userRole.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "invalid role type"})
			return
		}

		// Check if role is in allowed list
		roleAllowed := false
		for _, allowed := range allowedRoles {
			if role == allowed {
				roleAllowed = true
				break
			}
		}

		if !roleAllowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			return
		}

		c.Next()
	}
}

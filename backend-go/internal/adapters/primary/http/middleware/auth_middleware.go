package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/ports/services"
)

// AuthMiddleware handles JWT validation and user context
type AuthMiddleware struct {
	authService services.AuthService
}

func NewAuthMiddleware(a services.AuthService) *AuthMiddleware {
	return &AuthMiddleware{authService: a}
}

// RequireAuth verifies JWT token and sets user ID in context
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "no token provided"})
			return
		}

		// Extract token from Bearer scheme
		token := strings.TrimPrefix(header, "Bearer ")
		if token == header {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token format"})
			return
		}

		// Validate token and get user ID
		userID, err := m.authService.ValidateToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		// Set user ID in context for handlers
		c.Set("userID", userID)
		c.Next()
	}
}

// RequireRole ensures user has required role
func (m *AuthMiddleware) RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: implement role check using userID from context
		c.Next()
	}
}

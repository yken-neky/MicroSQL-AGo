package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/pkg/logging"
	"go.uber.org/zap"
)

// LoggingMiddleware attaches a request-scoped logger to the Gin context.
type LoggingMiddleware struct {
	Logger *zap.Logger
}

// NewLoggingMiddleware creates a new instance bound to the base logger.
func NewLoggingMiddleware(logger *zap.Logger) *LoggingMiddleware {
	return &LoggingMiddleware{Logger: logger}
}

// RequestLogger returns a gin.HandlerFunc that injects a logger with request_id and basic fields.
func (m *LoggingMiddleware) RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		reqID := uuid.New().String()

		// Build a child logger with request-specific fields
		rl := m.Logger.With(
			zap.String("request_id", reqID),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("remote_ip", c.ClientIP()),
		)

		// Redact sensitive headers (Authorization)
		auth := c.GetHeader("Authorization")
		if auth != "" {
			rl = rl.With(zap.String("authorization", logging.RedactAuthHeader(auth)))
		}

		// attach logger to context for handlers to reuse
		c.Set("logger", rl)

		rl.Info("request.start")

		// process request
		c.Next()

		// After handler
		duration := time.Since(start)
		status := c.Writer.Status()
		rl = rl.With(zap.Int("status", status), zap.Duration("duration_ms", duration))
		rl.Info("request.done")
	}
}

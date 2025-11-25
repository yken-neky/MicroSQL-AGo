package middleware

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

func TestRequestLoggerAttachesLogger(t *testing.T) {
    gin.SetMode(gin.TestMode)
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)

    req := httptest.NewRequest(http.MethodGet, "/", nil)
    req.Header.Set("Authorization", "Bearer token1234567890abcd")
    c.Request = req

    base := zap.NewNop()
    mw := NewLoggingMiddleware(base)

    handler := mw.RequestLogger()
    // call middleware
    handler(c)

    // logger should be present in context
    v, ok := c.Get("logger")
    if !ok {
        t.Fatal("expected logger to be set in context")
    }
    if _, ok := v.(*zap.Logger); !ok {
        t.Fatalf("expected *zap.Logger in context, got %T", v)
    }
}

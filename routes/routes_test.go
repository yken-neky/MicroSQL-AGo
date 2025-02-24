package routes_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"main/routes"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSetupRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	routes.SetupRoutes(router)

	tests := []struct {
		method       string
		url          string
		expectedCode int
	}{
		{"POST", "/api/v1/controles", http.StatusOK},
		{"GET", "/api/v1/controles", http.StatusOK},
		{"GET", "/api/v1/controles/1", http.StatusOK},
		{"PUT", "/api/v1/controles/1", http.StatusOK},
		{"DELETE", "/api/v1/controles/1", http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.method+" "+tt.url, func(t *testing.T) {
			req, _ := http.NewRequest(tt.method, tt.url, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}

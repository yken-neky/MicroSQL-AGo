package routes

import (
	"main/internal/controllers"

	"github.com/gin-gonic/gin"
)

// SetupRoutes initializes the API routes for the application.
// It sets up the routes for the /api/v1/controles endpoint, including
// creating, retrieving, updating, and deleting controls.
//
// Parameters:
//   - router: A pointer to the gin.Engine instance used to define the routes.
func SetupRoutes(router *gin.Engine) {
	cc := controllers.NewController()

	api := router.Group("/api/v1/controles")
	{
		api.POST("", cc.CreateControl)
		api.GET("", cc.GetControls)
		api.GET("/:id", cc.GetControlByID)
		api.PUT("/:id", cc.UpdateControl)
		api.DELETE("/:id", cc.DeleteControl)
	}
}

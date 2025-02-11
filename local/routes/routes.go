package routes

import (
	"main/local/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	cc := controllers.NewController()

	api := router.Group("/api/v1/controles")
	{
		api.POST("", cc.CreateControl)
		api.GET("", cc.GetControls)
		// api.GET("/:id", cc.GetControlByID)
		// api.PUT("/:id", cc.UpdateControl)
		// api.DELETE("/:id", cc.DeleteControl)
	}
}

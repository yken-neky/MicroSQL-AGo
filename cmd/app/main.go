package main

import (
	"main/routes"

	"github.com/gin-gonic/gin"
)

func main() {

	// Clinician server
	router := gin.Default()
	routes.SetupRoutes(router)

	err := router.Run(":8080")
	if err != nil {
		return
	}
}

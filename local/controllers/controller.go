package controllers

import (
	"net/http"

	"gorm.io/gorm"

	"main/local/connections"
	"main/local/models"

	"github.com/gin-gonic/gin"
)

type ControlController struct {
	DB      *gorm.DB
	ReqChan chan func() // Canal para manejar concurrencia
}

func NewController() *ControlController {
	cc := &ControlController{
		DB:      connections.DB,
		ReqChan: make(chan func()),
	}
	go cc.processRequests() // Iniciar worker
	return cc
}

// Worker para procesar solicitudes secuencialmente
func (cc *ControlController) processRequests() {
	for f := range cc.ReqChan {
		f()
	}
}

// CRUD Handlers con manejo de concurrencia
func (cc *ControlController) CreateControl(c *gin.Context) {
	var control models.Control
	if err := c.ShouldBindJSON(&control); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cc.ReqChan <- func() {
		if err := cc.DB.Create(&control).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, control)
	}
}

func (cc *ControlController) GetControls(c *gin.Context) {
	var controls []models.Control

	cc.ReqChan <- func() {
		if err := cc.DB.Find(&controls).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, controls)
	}
}

// ... Implementar similares para GetControlByID, UpdateControl, DeleteControl

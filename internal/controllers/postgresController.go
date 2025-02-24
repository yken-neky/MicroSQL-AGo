package controllers

import (
	"errors"
	"net/http"

	"gorm.io/gorm"

	"main/internal/connections"
	"main/internal/models"

	"github.com/gin-gonic/gin"
)

type ControlController struct {
	DB      *gorm.DB
	ReqChan chan func()
}

func NewController() *ControlController {
	cc := &ControlController{
		DB:      connections.SetupDatabase(),
		ReqChan: make(chan func()),
	}
	go cc.processRequests()
	return cc
}

func (cc *ControlController) processRequests() {
	for f := range cc.ReqChan {
		f()
	}
}

// CreateControl maneja la creación sincronizada
func (cc *ControlController) CreateControl(c *gin.Context) {
	var inputControl models.Control
	if err := c.ShouldBindJSON(&inputControl); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	type result struct {
		control models.Control
		err     error
	}
	resultChan := make(chan result)

	cc.ReqChan <- func() {
		controlToCreate := inputControl // Copia local para evitar race conditions
		err := cc.DB.Create(&controlToCreate).Error
		resultChan <- result{control: controlToCreate, err: err}
	}

	res := <-resultChan
	if res.err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": res.err.Error()})
		return
	}
	c.JSON(http.StatusCreated, res.control)
}

// GetControls obtiene todos los controles
func (cc *ControlController) GetControls(c *gin.Context) {
	type result struct {
		controls []models.Control
		err      error
	}
	resultChan := make(chan result)

	cc.ReqChan <- func() {
		var controls []models.Control
		err := cc.DB.Find(&controls).Error
		resultChan <- result{controls: controls, err: err}
	}

	res := <-resultChan
	if res.err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": res.err.Error()})
		return
	}
	c.JSON(http.StatusOK, res.controls)
}

// GetControlByID obtiene un control por ID
func (cc *ControlController) GetControlByID(c *gin.Context) {
	id := c.Param("id")

	type result struct {
		control models.Control
		err     error
	}
	resultChan := make(chan result)

	cc.ReqChan <- func() {
		var control models.Control
		err := cc.DB.First(&control, id).Error
		resultChan <- result{control: control, err: err}
	}

	res := <-resultChan
	if errors.Is(res.err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Control no encontrado"})
		return
	}
	if res.err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": res.err.Error()})
		return
	}
	c.JSON(http.StatusOK, res.control)
}

// UpdateControl actualiza un control
func (cc *ControlController) UpdateControl(c *gin.Context) {
	var inputControl models.Control
	if err := c.ShouldBindJSON(&inputControl); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	type result struct {
		control      models.Control
		rowsAffected int64
		err          error
	}
	resultChan := make(chan result)

	cc.ReqChan <- func() {
		controlToUpdate := inputControl
		dbResult := cc.DB.Save(&controlToUpdate)
		resultChan <- result{
			control:      controlToUpdate,
			rowsAffected: dbResult.RowsAffected,
			err:          dbResult.Error,
		}
	}

	res := <-resultChan
	if res.err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": res.err.Error()})
		return
	}
	if res.rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Control no encontrado"})
		return
	}
	c.JSON(http.StatusOK, res.control)
}

// DeleteControl elimina un control
func (cc *ControlController) DeleteControl(c *gin.Context) {
	id := c.Param("id")

	type result struct {
		rowsAffected int64
		err          error
	}
	resultChan := make(chan result)

	cc.ReqChan <- func() {
		dbResult := cc.DB.Delete(&models.Control{}, id)
		resultChan <- result{
			rowsAffected: dbResult.RowsAffected,
			err:          dbResult.Error,
		}
	}

	res := <-resultChan
	if res.err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": res.err.Error()})
		return
	}
	if res.rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Control no encontrado"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
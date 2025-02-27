package controllers

import (
	"errors"
	"gorm.io/gorm"
	"main/internal/utils"
	"net/http"

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
	var dto models.ControlDTO // Usar DTO en lugar del modelo directo
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resultChan := make(chan models.GetcOne)

	cc.ReqChan <- func() {
		control := models.Control{
			Nombre:      dto.Nombre,
			Descripcion: dto.Descripcion,
			Estado:      dto.Estado,
		}
		err := cc.DB.Create(&control).Error
		resultChan <- models.GetcOne{Control: control, Err: err}
	}

	res := <-resultChan
	if res.Err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": res.Err.Error()})
		return
	}
	c.JSON(http.StatusCreated, res.Control)
}

// GetControls obtiene todos los controles
func (cc *ControlController) GetControls(c *gin.Context) {

	resultChan := make(chan models.GetcAll)

	cc.ReqChan <- func() {
		var controls []models.Control
		err := cc.DB.Find(&controls).Error
		resultChan <- models.GetcAll{Controls: controls, Err: err}
	}

	res := <-resultChan
	if res.Err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": res.Err.Error()})
		return
	}
	c.JSON(http.StatusOK, res.Controls)
}

// GetControlByID obtiene un control por ID
func (cc *ControlController) GetControlByID(c *gin.Context) {
	id := c.Param("id")

	resultChan := make(chan models.GetcOne)

	cc.ReqChan <- func() {
		var control models.Control
		err := cc.DB.Where("id = ?", id).First(&control).Error // WHERE explícito
		resultChan <- models.GetcOne{Control: control, Err: err}
	}

	res := <-resultChan
	if errors.Is(res.Err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Control no encontrado"})
		return
	}
	if res.Err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": res.Err.Error()})
		return
	}
	c.JSON(http.StatusOK, res.Control)
}

// UpdateControl actualiza un control
func (cc *ControlController) UpdateControl(c *gin.Context) {
	id := c.Param("id") // Obtener ID del URL

	// 1. Validar formato del ID (ejemplo para UUID)
	if !utils.IsValidUUID(id) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// 2. Usar DTO para bindear solo campos actualizables

	var dto models.ControlDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validar campos requeridos
	// if dto.Nombre == "" {
	//  	c.JSON(http.StatusBadRequest, gin.H{"error": "Nombre es requerido"})
	//  	return
	// }

	resultChan := make(chan error)
	var updatedControl models.Control

	cc.ReqChan <- func() {
		//Copia local para evitar race condition
		localDTO := dto

		//Verificar existencia
		var existingControl models.Control
		if err := cc.DB.Where("id = ?", id).First(&existingControl).Error; err != nil {
			resultChan <- err
			return
		}

		//Actualizar solo campos permitidos
		updateData := map[string]interface{}{
			"nombre":      localDTO.Nombre,
			"descripcion": localDTO.Descripcion,
			"estado":      localDTO.Estado,
		}

		//Ejecutar actualización parcial
		err := cc.DB.Model(&existingControl).Updates(updateData).Error
		if err == nil {
			cc.DB.First(&existingControl, id) // Recargar datos en base de datos
		}
		updatedControl = existingControl
		resultChan <- err
	}

	if err := <-resultChan; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Control no encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedControl)
}

// DeleteControl elimina un control
func (cc *ControlController) DeleteControl(c *gin.Context) {
	id := c.Param("id")

	// Validar formato del ID (ejemplo para UUID)
	if !utils.IsValidUUID(id) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	resultChan := make(chan error)

	cc.ReqChan <- func() {
		// Verificar existencia del registro
		var control models.Control
		if err := cc.DB.Where("id = ?", id).First(&control).Error; err != nil {
			resultChan <- err // Envía error si no existe
			return
		}

		// Borrado físico (opcional: usar soft-delete si está configurado)
		if err := cc.DB.Where("id = ?", id).Delete(&control).Error; err != nil {
			resultChan <- err
			return
		}

		resultChan <- nil
	}

	err := <-resultChan
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Control no encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 4. Respuesta correcta sin cuerpo
	c.Status(http.StatusNoContent)
}

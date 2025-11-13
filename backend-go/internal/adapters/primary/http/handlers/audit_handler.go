package handlers

import (
	"net/http"
	
	"github.com/gin-gonic/gin"
	controlsuc "github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/usecases/controls"
)

// AuditHandler maneja endpoints relacionados con auditorías de controles
type AuditHandler struct {
	auditUC *controlsuc.ExecuteAuditUseCase
}

func NewAuditHandler(a *controlsuc.ExecuteAuditUseCase) *AuditHandler {
	return &AuditHandler{auditUC: a}
}

// ExecuteAudit ejecuta una auditoría parcial o completa según request
func (h *AuditHandler) ExecuteAudit(c *gin.Context) {
	var req controlsuc.AuditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userID")

	res, err := h.auditUC.Execute(c.Request.Context(), userID.(uint), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

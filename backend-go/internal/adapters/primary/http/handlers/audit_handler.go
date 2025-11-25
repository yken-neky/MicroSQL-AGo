package handlers

import (
	"fmt"
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

	// return audit result with audit_run_id
	c.JSON(http.StatusOK, res)
}

// GetAudit returns details for a specific audit run
func (h *AuditHandler) GetAudit(c *gin.Context) {
	idParam := c.Param("id")
	var id uint64
	var err error
	if id, err = parseUintParam(idParam); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid audit id"})
		return
	}

	userID, _ := c.Get("userID")

	res, run, err := h.auditUC.GetAuditRun(c.Request.Context(), userID.(uint), uint(id))
	if err != nil {
		if err.Error() == "forbidden" {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"audit": run, "result": res})
}

// parseUintParam converts string id parameter to uint64
func parseUintParam(s string) (uint64, error) {
	var id uint64
	_, err := fmt.Sscanf(s, "%d", &id)
	return id, err
}

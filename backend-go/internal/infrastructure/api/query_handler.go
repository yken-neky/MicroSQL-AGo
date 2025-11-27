package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/ports/repositories"
	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/usecases/queries"
)

// QueryHandler maneja las peticiones HTTP relacionadas con consultas
type QueryHandler struct {
	executeQuery *queries.ExecuteQueryUseCase
	queryRepo    repositories.QueryRepository
}

func NewQueryHandler(eq *queries.ExecuteQueryUseCase, qr repositories.QueryRepository) *QueryHandler {
	return &QueryHandler{
		executeQuery: eq,
		queryRepo:    qr,
	}
}

// ExecuteQuery maneja la ejecución de consultas SQL
func (h *QueryHandler) ExecuteQuery(c *gin.Context) {
	// For security, arbitrary SQL execution via the API is disallowed.
	// Users must execute predefined control scripts via the audits endpoint (now per-manager): /api/db/{gestor}/audits/execute
	c.JSON(http.StatusForbidden, gin.H{"error": "execution of raw SQL via API is forbidden; use /api/db/{gestor}/audits/execute to run predefined control scripts"})
}

// GetQueryResult obtiene el resultado de una consulta por ID
func (h *QueryHandler) GetQueryResult(c *gin.Context) {
	queryID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "100"))

	result, err := h.queryRepo.GetResult(uint(queryID), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetQueryStats obtiene las estadísticas de una consulta
func (h *QueryHandler) GetQueryStats(c *gin.Context) {
	queryID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query ID"})
		return
	}

	stats, err := h.queryRepo.GetStats(uint(queryID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// ListUserQueries lista las consultas de un usuario
func (h *QueryHandler) ListUserQueries(c *gin.Context) {
	userID, _ := c.Get("userID")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	queries, err := h.queryRepo.ListByUser(userID.(uint), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, queries)
}

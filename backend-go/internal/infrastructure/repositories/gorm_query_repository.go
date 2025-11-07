package repositories

import (
	"encoding/json"
	"time"

	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/entities"
	"gorm.io/gorm"
)

// QueryResultDB is the DB representation for saved query results
type QueryResultDB struct {
	QueryID     uint   `gorm:"primaryKey"`
	Columns     string `gorm:"type:text"` // JSON de columnas
	Types       string `gorm:"type:text"` // JSON de tipos
	Rows        string `gorm:"type:text"` // JSON de filas
	HasMoreRows bool
	PageSize    int
	PageNumber  int
	CreatedAt   time.Time
}

// GormQueryRepository implementa QueryRepository usando GORM
type GormQueryRepository struct {
	db *gorm.DB
}

func NewGormQueryRepository(db *gorm.DB) *GormQueryRepository {
	return &GormQueryRepository{db: db}
}

// Create crea un nuevo registro de consulta
func (r *GormQueryRepository) Create(query *entities.Query) error {
	return r.db.Create(query).Error
}

// Update actualiza un registro de consulta existente
func (r *GormQueryRepository) Update(query *entities.Query) error {
	return r.db.Save(query).Error
}

// GetByID obtiene una consulta por su ID
func (r *GormQueryRepository) GetByID(id uint) (*entities.Query, error) {
	var query entities.Query
	if err := r.db.First(&query, id).Error; err != nil {
		return nil, err
	}
	return &query, nil
}

// ListByUser lista las consultas de un usuario con paginación
func (r *GormQueryRepository) ListByUser(userID uint, page, pageSize int) ([]entities.Query, error) {
	var queries []entities.Query
	offset := (page - 1) * pageSize

	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&queries).Error

	if err != nil {
		return nil, err
	}

	return queries, nil
}

// SaveResult guarda el resultado de una consulta SELECT
func (r *GormQueryRepository) SaveResult(result *entities.QueryResult) error {
	// Convertir arrays a JSON
	columnsJSON, err := json.Marshal(result.Columns)
	if err != nil {
		return err
	}

	typesJSON, err := json.Marshal(result.Types)
	if err != nil {
		return err
	}

	rowsJSON, err := json.Marshal(result.Rows)
	if err != nil {
		return err
	}

	resultDB := QueryResultDB{
		QueryID:     result.QueryID,
		Columns:     string(columnsJSON),
		Types:       string(typesJSON),
		Rows:        string(rowsJSON),
		HasMoreRows: result.HasMoreRows,
		PageSize:    result.PageSize,
		PageNumber:  result.PageNumber,
		CreatedAt:   time.Now(),
	}

	return r.db.Create(&resultDB).Error
}

// GetResult obtiene el resultado de una consulta por su ID
func (r *GormQueryRepository) GetResult(queryID uint, page, pageSize int) (*entities.QueryResult, error) {
	var resultDB QueryResultDB
	if err := r.db.Where("query_id = ?", queryID).First(&resultDB).Error; err != nil {
		return nil, err
	}

	var columns []string
	var types []string
	var rows [][]any

	if err := json.Unmarshal([]byte(resultDB.Columns), &columns); err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(resultDB.Types), &types); err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(resultDB.Rows), &rows); err != nil {
		return nil, err
	}

	return &entities.QueryResult{
		QueryID:     resultDB.QueryID,
		Columns:     columns,
		Types:       types,
		Rows:        rows,
		HasMoreRows: resultDB.HasMoreRows,
		PageSize:    resultDB.PageSize,
		PageNumber:  resultDB.PageNumber,
	}, nil
}

// SaveStats guarda las estadísticas de ejecución
func (r *GormQueryRepository) SaveStats(stats *entities.ExecutionStats) error {
	return r.db.Create(stats).Error
}

// GetStats obtiene las estadísticas de ejecución por ID de consulta
func (r *GormQueryRepository) GetStats(queryID uint) (*entities.ExecutionStats, error) {
	var stats entities.ExecutionStats
	if err := r.db.Where("query_id = ?", queryID).First(&stats).Error; err != nil {
		return nil, err
	}
	return &stats, nil
}

// GetUserQueryCount obtiene el total de consultas de un usuario
func (r *GormQueryRepository) GetUserQueryCount(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&entities.Query{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

// CleanOldQueries elimina consultas antiguas
func (r *GormQueryRepository) CleanOldQueries(olderThan string) error {
	return r.db.Where("created_at < ?", olderThan).Delete(&entities.Query{}).Error
}

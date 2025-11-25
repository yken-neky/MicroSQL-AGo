package repositories

import (
    "github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/entities"
    "gorm.io/gorm"
)

// GormAuditRepository implements AuditRepository using GORM
type GormAuditRepository struct {
    db *gorm.DB
}

func NewGormAuditRepository(db *gorm.DB) *GormAuditRepository {
    return &GormAuditRepository{db: db}
}

func (r *GormAuditRepository) CreateAuditRun(run *entities.AuditRun) error {
    return r.db.Create(run).Error
}

func (r *GormAuditRepository) UpdateAuditRun(run *entities.AuditRun) error {
    return r.db.Save(run).Error
}

func (r *GormAuditRepository) GetAuditRunByID(id uint) (*entities.AuditRun, error) {
    var run entities.AuditRun
    if err := r.db.First(&run, id).Error; err != nil {
        return nil, err
    }
    return &run, nil
}

func (r *GormAuditRepository) CreateScriptResult(res *entities.AuditScriptResult) error {
    return r.db.Create(res).Error
}

func (r *GormAuditRepository) ListScriptResultsByAuditRun(auditRunID uint) ([]entities.AuditScriptResult, error) {
    var list []entities.AuditScriptResult
    if err := r.db.Where("audit_run_id = ?", auditRunID).Find(&list).Error; err != nil {
        return nil, err
    }
    return list, nil
}

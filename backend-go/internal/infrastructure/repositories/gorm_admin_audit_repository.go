package repositories

import (
	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/entities"
	"gorm.io/gorm"
)

type GormAdminAuditRepository struct {
	db *gorm.DB
}

func NewGormAdminAuditRepository(db *gorm.DB) *GormAdminAuditRepository {
	return &GormAdminAuditRepository{db: db}
}

func (r *GormAdminAuditRepository) Create(log *entities.AdminActionLog) error {
	return r.db.Create(log).Error
}

func (r *GormAdminAuditRepository) List(actorID *uint, targetType *string, action *string, limit int, offset int) ([]entities.AdminActionLog, error) {
	var list []entities.AdminActionLog
	q := r.db.Model(&entities.AdminActionLog{})
	if actorID != nil {
		q = q.Where("actor_id = ?", *actorID)
	}
	if targetType != nil {
		q = q.Where("target_type = ?", *targetType)
	}
	if action != nil {
		q = q.Where("action = ?", *action)
	}
	if limit <= 0 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	if err := q.Order("created_at DESC").Limit(limit).Offset(offset).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

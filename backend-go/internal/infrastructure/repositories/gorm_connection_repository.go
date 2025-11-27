package repositories

import (
	"time"

	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/entities"
	"gorm.io/gorm"
)

// GormConnectionRepository implementa ConnectionRepository usando GORM
type GormConnectionRepository struct {
	db *gorm.DB
}

func NewGormConnectionRepository(db *gorm.DB) *GormConnectionRepository {
	return &GormConnectionRepository{db: db}
}

func (r *GormConnectionRepository) CreateActive(conn *entities.ActiveConnection) error {
	return r.db.Create(conn).Error
}

func (r *GormConnectionRepository) GetActiveByUserID(userID uint) (*entities.ActiveConnection, error) {
	var ac entities.ActiveConnection
	if err := r.db.Where("user_id = ?", userID).First(&ac).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &ac, nil
}

func (r *GormConnectionRepository) GetActiveByUserIDAndDriver(userID uint, driver string) (*entities.ActiveConnection, error) {
	var ac entities.ActiveConnection
	if err := r.db.Where("user_id = ? AND driver = ? AND is_connected = ?", userID, driver, true).First(&ac).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &ac, nil
}

func (r *GormConnectionRepository) UpdateActive(conn *entities.ActiveConnection) error {
	return r.db.Save(conn).Error
}

func (r *GormConnectionRepository) DeleteActive(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&entities.ActiveConnection{}).Error
}

func (r *GormConnectionRepository) ListActive() ([]*entities.ActiveConnection, error) {
	var list []*entities.ActiveConnection
	if err := r.db.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *GormConnectionRepository) ListActiveByUser(userID uint) ([]*entities.ActiveConnection, error) {
	var list []*entities.ActiveConnection
	if err := r.db.Where("user_id = ? AND is_connected = ?", userID, true).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *GormConnectionRepository) LogConnection(log *entities.ConnectionLog) error {
	log.Timestamp = time.Now()
	return r.db.Create(log).Error
}

func (r *GormConnectionRepository) GetLogsByUserID(userID uint, limit, offset int) ([]*entities.ConnectionLog, error) {
	var logs []*entities.ConnectionLog
	q := r.db.Where("user_id = ?", userID).Order("created_at DESC").Offset(offset)
	if limit > 0 {
		q = q.Limit(limit)
	}
	if err := q.Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}

func (r *GormConnectionRepository) GetLogByID(id uint) (*entities.ConnectionLog, error) {
	var log entities.ConnectionLog
	if err := r.db.First(&log, id).Error; err != nil {
		return nil, err
	}
	return &log, nil
}

func (r *GormConnectionRepository) CountLogsByUserID(userID uint) (int64, error) {
	var cnt int64
	if err := r.db.Model(&entities.ConnectionLog{}).Where("user_id = ?", userID).Count(&cnt).Error; err != nil {
		return 0, err
	}
	return cnt, nil
}

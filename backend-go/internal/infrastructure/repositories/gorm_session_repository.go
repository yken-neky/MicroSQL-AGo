package repositories

import (
	"time"

	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/entities"
	"gorm.io/gorm"
)

// GormSessionRepository persists sessions
type GormSessionRepository struct {
	db *gorm.DB
}

func NewGormSessionRepository(db *gorm.DB) *GormSessionRepository {
	return &GormSessionRepository{db: db}
}

func (r *GormSessionRepository) CreateSession(s *entities.Session) error {
	return r.db.Create(s).Error
}

func (r *GormSessionRepository) GetActiveByUserID(userID uint) (*entities.Session, error) {
	var s entities.Session
	if err := r.db.Where("user_id = ? AND is_active = ?", userID, true).First(&s).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	// if expired, mark inactive and return nil
	if s.ExpiresAt != nil && s.ExpiresAt.Before(time.Now()) {
		s.IsActive = false
		_ = r.db.Save(&s).Error
		return nil, nil
	}
	return &s, nil
}

func (r *GormSessionRepository) GetByToken(token string) (*entities.Session, error) {
	var s entities.Session
	if err := r.db.Where("token = ?", token).First(&s).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *GormSessionRepository) DeactivateByToken(token string) error {
	return r.db.Model(&entities.Session{}).Where("token = ?", token).Update("is_active", false).Error
}

// ListActiveSessions returns sessions where is_active is true and not yet expired
func (r *GormSessionRepository) ListActiveSessions() ([]entities.Session, error) {
	var list []entities.Session
	now := time.Now()
	q := r.db.Where("is_active = ?", true)
	// also treat sessions with ExpiresAt in the past as inactive (skip them)
	q = q.Where("expires_at IS NULL OR expires_at > ?", now)
	if err := q.Order("created_at DESC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

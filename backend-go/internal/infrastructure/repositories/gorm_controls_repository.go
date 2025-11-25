package repositories

import (
	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/entities"
	repoport "github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/ports/repositories"
	"gorm.io/gorm"
)

// GormControlsRepository implementa ControlRepository usando GORM
type GormControlsRepository struct {
	db *gorm.DB
}

func NewGormControlsRepository(db *gorm.DB) *GormControlsRepository {
	return &GormControlsRepository{db: db}
}

// ListControls lista todos los controles
func (r *GormControlsRepository) ListControls() ([]entities.ControlsInformation, error) {
	var controls []entities.ControlsInformation
	if err := r.db.Find(&controls).Error; err != nil {
		return nil, err
	}
	return controls, nil
}

// GetControlScripts obtiene los scripts asociados a un control
func (r *GormControlsRepository) GetControlScripts(controlID uint) ([]repoport.ControlsScript, error) {
	var scripts []repoport.ControlsScript
	if err := r.db.Where("control_script_id = ?", controlID).Find(&scripts).Error; err != nil {
		return nil, err
	}
	return scripts, nil
}

// GetScriptsByIDs obtiene scripts por su ID
func (r *GormControlsRepository) GetScriptsByIDs(ids []uint) ([]repoport.ControlsScript, error) {
	var scripts []repoport.ControlsScript
	if err := r.db.Where("id IN ?", ids).Find(&scripts).Error; err != nil {
		return nil, err
	}
	return scripts, nil
}

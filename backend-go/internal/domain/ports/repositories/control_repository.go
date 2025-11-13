package repositories

import "github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/entities"

// ControlRepository define operaciones para leer controles y sus scripts
type ControlRepository interface {
	// ListControls lista todos los controles disponibles
	ListControls() ([]entities.ControlsInformation, error)

	// GetControlScripts obtiene los scripts asociados a un control
	GetControlScripts(controlID uint) ([]ControlsScript, error)

	// GetScriptsByIDs obtiene scripts por su ID
	GetScriptsByIDs(ids []uint) ([]ControlsScript, error)
}

// ControlsScript es la representación en repositorio de un script de control
type ControlsScript struct {
	ID               uint   `gorm:"primaryKey" json:"id"`
	ControlType      string `gorm:"column:control_type" json:"control_type"`
	QuerySQL         string `gorm:"column:query_sql;type:text" json:"query_sql"`
	ControlScriptRef uint   `gorm:"column:control_script_id" json:"control_id"`
}

func (ControlsScript) TableName() string { return "controls_scripts" }

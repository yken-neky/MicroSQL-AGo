package entities

// ControlsScript representa un script de control asociado a un control
type ControlsScript struct {
    ID          uint   `gorm:"primaryKey" json:"id"`
    ControlType string `gorm:"column:control_type" json:"control_type"`
    QuerySQL    string `gorm:"column:query_sql;type:text" json:"query_sql"`
    ControlID   uint   `gorm:"column:control_script_id_id" json:"control_id"`
}

func (ControlsScript) TableName() string { return "controls_scripts" }

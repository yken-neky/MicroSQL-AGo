package models

import "gorm.io/gorm"

type Control struct {
	gorm.Model
	Nombre      string `gorm:"not null" json:"nombre"`
	Descripcion string `json:"descripcion"`
	Estado      bool   `json:"estado"`
}

type ControlDTO struct {
	Nombre      string `json:"nombre" binding:"required"`
	Descripcion string `json:"descripcion" binding:"required"`
	Estado      bool   `json:"estado"`
}

type GetcOne struct {
	Control Control
	Err     error
}

type GetcAll struct {
	Controls []Control
	Err      error
}

// Migrar la tabla
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&Control{})
}

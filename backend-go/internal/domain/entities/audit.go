package entities

import "time"

// AuditRun represents a single audit execution (batch of control scripts)
type AuditRun struct {
    ID         uint       `gorm:"primaryKey" json:"id"`
    UserID     uint       `gorm:"not null;index" json:"user_id"`
    Mode       string     `gorm:"size:20;not null;default:'partial'" json:"mode"` // partial|full
    Database   string     `gorm:"size:255" json:"database"`
    Total      int        `json:"total"`
    Passed     int        `json:"passed"`
    Failed     int        `json:"failed"`
    Status     string     `gorm:"size:20;index;default:'running'" json:"status"`
    Controls   string     `gorm:"type:text" json:"controls"` // JSON array of control IDs (optional)
    StartedAt  time.Time  `gorm:"autoCreateTime" json:"started_at"`
    FinishedAt *time.Time `gorm:"column:finished_at" json:"finished_at"`
}

// AuditScriptResult represents result of executing one control script inside an audit run.
type AuditScriptResult struct {
    ID         uint      `gorm:"primaryKey" json:"id"`
    AuditRunID uint      `gorm:"index;not null" json:"audit_run_id"`
    ScriptID   uint      `gorm:"not null;index" json:"script_id"`
    ControlID  uint      `gorm:"not null;index" json:"control_id"`
    QuerySQL   string    `gorm:"type:text" json:"query_sql"`
    Passed     bool      `json:"passed"`
    Error      string    `gorm:"type:text" json:"error"`
    DurationMs int64     `json:"duration_ms"`
    Rows       int64     `json:"rows"`
    CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
}

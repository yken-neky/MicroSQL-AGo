package repositories

import "github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/entities"

// AdminAuditRepository persists AdminActionLog records
type AdminAuditRepository interface {
	Create(log *entities.AdminActionLog) error
	// List returns audit logs matching optional filters (actorID, targetType, action)
	List(actorID *uint, targetType *string, action *string, limit int, offset int) ([]entities.AdminActionLog, error)
}

package repositories

import "github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/entities"

// SessionRepository defines operations to manage user sessions
type SessionRepository interface {
	CreateSession(s *entities.Session) error
	GetActiveByUserID(userID uint) (*entities.Session, error)
	GetByToken(token string) (*entities.Session, error)
	DeactivateByToken(token string) error
	// ListActiveSessions returns all currently active (non-deactivated, non-expired) sessions
	ListActiveSessions() ([]entities.Session, error)
}

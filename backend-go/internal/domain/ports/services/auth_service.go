package services

import "github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/entities"

// AuthService defines methods for authentication and token management
type AuthService interface {
	HashPassword(plain string) (string, error)
	CheckPassword(hash, plain string) error
	GenerateToken(user *entities.User) (string, error)
	ValidateToken(token string) (uint, error) // returns user ID
}

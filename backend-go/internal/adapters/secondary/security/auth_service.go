package security

import (
	"errors"
	"time"

	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/adapters/secondary/encryption"
	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/entities"
)

// Validation errors
var (
	ErrInvalidToken = errors.New("token is invalid or expired")
	ErrInvalidCreds = errors.New("invalid credentials")
)

// AuthService implements ports.AuthService using JWT
type AuthService struct {
	jwtSvc   *JWTService
	expHours int
}

func NewAuthService(secret string, expirationHours int) *AuthService {
	return &AuthService{
		jwtSvc:   NewJWTService(secret),
		expHours: expirationHours,
	}
}

func (s *AuthService) HashPassword(plain string) (string, error) {
	return encryption.HashPassword(plain)
}

func (s *AuthService) CheckPassword(hash, plain string) error {
	return encryption.CheckPassword(hash, plain)
}

func (s *AuthService) GenerateToken(user *entities.User) (string, error) {
	if user == nil || user.ID == 0 {
		return "", errors.New("invalid user")
	}

	// Update last login
	user.LastLogin = time.Now()

	return s.jwtSvc.Generate(user.ID, s.expHours)
}

func (s *AuthService) ValidateToken(token string) (uint, error) {
	if token == "" {
		return 0, ErrInvalidToken
	}
	return s.jwtSvc.Validate(token)
}

package user

import (
	"errors"

	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/entities"
	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/ports/repositories"
	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/ports/services"
)

// LoginUseCase handles user authentication
type LoginUseCase struct {
	userRepo    repositories.UserRepository
	authService services.AuthService
}

func NewLoginUseCase(r repositories.UserRepository, a services.AuthService) *LoginUseCase {
	return &LoginUseCase{userRepo: r, authService: a}
}

// Execute authenticates a user and returns JWT token
func (uc *LoginUseCase) Execute(username, password string) (*entities.User, string, error) {
	if username == "" || password == "" {
		return nil, "", errors.New("missing credentials")
	}

	user, err := uc.userRepo.FindByUsername(username)
	if err != nil {
		return nil, "", err
	}
	if !user.IsActive {
		return nil, "", errors.New("user is inactive")
	}

	if err := uc.authService.CheckPassword(user.Password, password); err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	token, err := uc.authService.GenerateToken(user)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

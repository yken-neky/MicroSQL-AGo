package user

import (
	"errors"

	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/entities"
	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/ports/repositories"
	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/ports/services"
)

// RegisterUserUseCase handles user registration
type RegisterUserUseCase struct {
	userRepo    repositories.UserRepository
	authService services.AuthService
}

func NewRegisterUserUseCase(r repositories.UserRepository, a services.AuthService) *RegisterUserUseCase {
	return &RegisterUserUseCase{userRepo: r, authService: a}
}

// Execute registers a new user and returns created user and JWT token
func (uc *RegisterUserUseCase) Execute(username, email, password string) (*entities.User, string, error) {
	if username == "" || email == "" || password == "" {
		return nil, "", errors.New("missing fields")
	}

	// check unique constraints at repo level
	if u, _ := uc.userRepo.FindByUsername(username); u != nil {
		return nil, "", errors.New("username already exists")
	}
	if e, _ := uc.userRepo.FindByEmail(email); e != nil {
		return nil, "", errors.New("email already exists")
	}

	hashed, err := uc.authService.HashPassword(password)
	if err != nil {
		return nil, "", err
	}

	user := &entities.User{
		Username: username,
		Email:    email,
		Password: hashed,
		Role:     "cliente",
		IsActive: true,
	}

	if err := uc.userRepo.Create(user); err != nil {
		return nil, "", err
	}

	token, err := uc.authService.GenerateToken(user)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

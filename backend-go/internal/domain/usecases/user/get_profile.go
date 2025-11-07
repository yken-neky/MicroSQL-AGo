package user

import (
	"errors"

	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/entities"
	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/ports/repositories"
)

// GetProfileUseCase retrieves user profile data
type GetProfileUseCase struct {
	userRepo repositories.UserRepository
}

func NewGetProfileUseCase(r repositories.UserRepository) *GetProfileUseCase {
	return &GetProfileUseCase{userRepo: r}
}

// Execute retrieves a user's profile by ID
func (uc *GetProfileUseCase) Execute(userID uint) (*entities.User, error) {
	if userID == 0 {
		return nil, errors.New("invalid user ID")
	}

	user, err := uc.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if !user.IsActive {
		return nil, errors.New("user is inactive")
	}

	return user, nil
}

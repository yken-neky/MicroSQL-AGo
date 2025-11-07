package roles

import (
	"context"
	"errors"

	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/ports/repositories"
)

// AssignRoleUseCase maneja la asignación de roles a usuarios
type AssignRoleUseCase struct {
	roleRepo repositories.RoleRepository
}

func NewAssignRoleUseCase(rr repositories.RoleRepository) *AssignRoleUseCase {
	return &AssignRoleUseCase{
		roleRepo: rr,
	}
}

func (uc *AssignRoleUseCase) Execute(ctx context.Context, userID uint, roleName string) error {
	// Obtener el rol por nombre
	role, err := uc.roleRepo.GetByName(roleName)
	if err != nil {
		return err
	}
	if role == nil {
		return errors.New("role not found")
	}

	// Asignar el rol al usuario
	return uc.roleRepo.AssignToUser(userID, role.ID)
}

// RevokeRoleUseCase maneja la revocación de roles de usuarios
type RevokeRoleUseCase struct {
	roleRepo repositories.RoleRepository
}

func NewRevokeRoleUseCase(rr repositories.RoleRepository) *RevokeRoleUseCase {
	return &RevokeRoleUseCase{
		roleRepo: rr,
	}
}

func (uc *RevokeRoleUseCase) Execute(ctx context.Context, userID uint, roleName string) error {
	// Obtener el rol por nombre
	role, err := uc.roleRepo.GetByName(roleName)
	if err != nil {
		return err
	}
	if role == nil {
		return errors.New("role not found")
	}

	// Revocar el rol del usuario
	return uc.roleRepo.RevokeFromUser(userID, role.ID)
}

package permissions

import (
	"context"
	"errors"

	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/ports/repositories"
)

// AssignPermissionUseCase maneja la asignación de permisos a roles
type AssignPermissionUseCase struct {
	permRepo repositories.PermissionRepository
	roleRepo repositories.RoleRepository
}

func NewAssignPermissionUseCase(
	pr repositories.PermissionRepository,
	rr repositories.RoleRepository,
) *AssignPermissionUseCase {
	return &AssignPermissionUseCase{
		permRepo: pr,
		roleRepo: rr,
	}
}

func (uc *AssignPermissionUseCase) Execute(ctx context.Context, roleName string, permissionName string) error {
	// Obtener el rol
	role, err := uc.roleRepo.GetByName(roleName)
	if err != nil {
		return err
	}
	if role == nil {
		return errors.New("role not found")
	}

	// Obtener el permiso
	perm, err := uc.permRepo.GetByName(permissionName)
	if err != nil {
		return err
	}
	if perm == nil {
		return errors.New("permission not found")
	}

	// Asignar el permiso al rol
	return uc.permRepo.AssignToRole(role.ID, perm.ID)
}

// RevokePermissionUseCase maneja la revocación de permisos de roles
type RevokePermissionUseCase struct {
	permRepo repositories.PermissionRepository
	roleRepo repositories.RoleRepository
}

func NewRevokePermissionUseCase(
	pr repositories.PermissionRepository,
	rr repositories.RoleRepository,
) *RevokePermissionUseCase {
	return &RevokePermissionUseCase{
		permRepo: pr,
		roleRepo: rr,
	}
}

func (uc *RevokePermissionUseCase) Execute(ctx context.Context, roleName string, permissionName string) error {
	// Obtener el rol
	role, err := uc.roleRepo.GetByName(roleName)
	if err != nil {
		return err
	}
	if role == nil {
		return errors.New("role not found")
	}

	// Obtener el permiso
	perm, err := uc.permRepo.GetByName(permissionName)
	if err != nil {
		return err
	}
	if perm == nil {
		return errors.New("permission not found")
	}

	// Revocar el permiso del rol
	return uc.permRepo.RevokeFromRole(role.ID, perm.ID)
}

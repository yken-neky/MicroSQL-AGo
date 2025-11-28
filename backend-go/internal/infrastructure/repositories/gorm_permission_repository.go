package repositories

import (
	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/entities"
	"gorm.io/gorm"
)

// GormPermissionRepository persists permissions and manages role-permission relations
type GormPermissionRepository struct {
	db *gorm.DB
}

func NewGormPermissionRepository(db *gorm.DB) *GormPermissionRepository {
	return &GormPermissionRepository{db: db}
}

func (r *GormPermissionRepository) Create(permission *entities.Permission) error {
	return r.db.Create(permission).Error
}

func (r *GormPermissionRepository) Update(permission *entities.Permission) error {
	return r.db.Save(permission).Error
}

func (r *GormPermissionRepository) Delete(id uint) error {
	return r.db.Delete(&entities.Permission{}, id).Error
}

func (r *GormPermissionRepository) GetByID(id uint) (*entities.Permission, error) {
	var p entities.Permission
	if err := r.db.First(&p, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func (r *GormPermissionRepository) GetByName(name string) (*entities.Permission, error) {
	var p entities.Permission
	if err := r.db.Where("name = ?", name).First(&p).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func (r *GormPermissionRepository) List() ([]entities.Permission, error) {
	var list []entities.Permission
	if err := r.db.Order("resource ASC, action ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *GormPermissionRepository) AssignToRole(roleID uint, permissionID uint) error {
	// Use GORM association API to attach permission to role
	role := entities.Role{ID: roleID}
	perm := entities.Permission{ID: permissionID}
	return r.db.Model(&role).Association("Permissions").Append(&perm)
}

func (r *GormPermissionRepository) RevokeFromRole(roleID uint, permissionID uint) error {
	role := entities.Role{ID: roleID}
	perm := entities.Permission{ID: permissionID}
	return r.db.Model(&role).Association("Permissions").Delete(&perm)
}

func (r *GormPermissionRepository) GetRolePermissions(roleID uint) ([]entities.Permission, error) {
	var role entities.Role
	if err := r.db.Preload("Permissions").First(&role, roleID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return role.Permissions, nil
}

func (r *GormPermissionRepository) HasPermission(roleID uint, permissionName string) (bool, error) {
	var count int64
	err := r.db.Table("role_permissions").Joins("JOIN permissions p ON p.id = role_permissions.permission_id").Where("role_permissions.role_id = ? AND p.name = ?", roleID, permissionName).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

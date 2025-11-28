package repositories

import (
	"github.com/yken-neky/MicroSQL-AGo/backend-go/internal/domain/entities"
	"gorm.io/gorm"
)

// GormRoleRepository persists roles and user-role relations
type GormRoleRepository struct {
	db *gorm.DB
}

func NewGormRoleRepository(db *gorm.DB) *GormRoleRepository {
	return &GormRoleRepository{db: db}
}

func (r *GormRoleRepository) Create(role *entities.Role) error {
	return r.db.Create(role).Error
}

func (r *GormRoleRepository) Update(role *entities.Role) error {
	return r.db.Save(role).Error
}

func (r *GormRoleRepository) Delete(id uint) error {
	return r.db.Delete(&entities.Role{}, id).Error
}

func (r *GormRoleRepository) GetByID(id uint) (*entities.Role, error) {
	var role entities.Role
	if err := r.db.Preload("Permissions").First(&role, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

func (r *GormRoleRepository) GetByName(name string) (*entities.Role, error) {
	var role entities.Role
	if err := r.db.Preload("Permissions").Where("name = ?", name).First(&role).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

func (r *GormRoleRepository) List() ([]entities.Role, error) {
	var roles []entities.Role
	if err := r.db.Preload("Permissions").Order("name ASC").Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *GormRoleRepository) AssignToUser(userID uint, roleID uint) error {
	ur := entities.UserRole{UserID: userID, RoleID: roleID}
	return r.db.Create(&ur).Error
}

func (r *GormRoleRepository) RevokeFromUser(userID uint, roleID uint) error {
	return r.db.Where("user_id = ? AND role_id = ?", userID, roleID).Delete(&entities.UserRole{}).Error
}

func (r *GormRoleRepository) GetUserRoles(userID uint) ([]entities.Role, error) {
	var roles []entities.Role
	// Join via user_roles table
	if err := r.db.Joins("JOIN user_roles ur ON ur.role_id = roles.id").Where("ur.user_id = ?", userID).Preload("Permissions").Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

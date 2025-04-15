package repository

import (
	"user_crud/internal/domain/entity"
	"user_crud/internal/domain/repository/interfaces"

	"gorm.io/gorm"
)

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) interfaces.RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) Create(role *entity.Role) error {
	return r.db.Create(role).Error
}

func (r *roleRepository) FindByName(name string) (entity.Role, error) {
	var role entity.Role
	err := r.db.Where("name = ?", name).First(&role).Error
	return role, err
}

func (r *roleRepository) FindByID(id uint) (entity.Role, error) {
	var role entity.Role
	err := r.db.First(&role, id).Error
	return role, err
}

func (r *roleRepository) FindAll() ([]entity.Role, error) {
	var roles []entity.Role
	err := r.db.Find(&roles).Error
	return roles, err
}

package interfaces

import (
	"user_crud/internal/domain/entity"
)

type RoleRepository interface {
	Create(role *entity.Role) error
	FindByName(name string) (entity.Role, error)
	FindByID(id uint) (entity.Role, error)
	FindAll() ([]entity.Role, error)
}

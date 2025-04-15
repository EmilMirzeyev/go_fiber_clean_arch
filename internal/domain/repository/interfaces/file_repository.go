package interfaces

import (
	"user_crud/internal/domain/entity"
)

type FileRepository interface {
	Create(file *entity.File) error
	FindByUserID(userID uint) (entity.File, error)
	Update(file *entity.File) error
	DeleteByUserID(userID uint) error
}

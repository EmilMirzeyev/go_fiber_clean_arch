package repository

import (
	"user_crud/internal/domain/entity"
	"user_crud/internal/domain/repository/interfaces"

	"gorm.io/gorm"
)

type fileRepository struct {
	db *gorm.DB
}

func NewFileRepository(db *gorm.DB) interfaces.FileRepository {
	return &fileRepository{db: db}
}

func (r *fileRepository) Create(file *entity.File) error {
	return r.db.Create(file).Error
}

func (r *fileRepository) FindByUserID(userID uint) (entity.File, error) {
	var file entity.File
	err := r.db.Where("user_id = ?", userID).First(&file).Error
	return file, err
}

func (r *fileRepository) Update(file *entity.File) error {
	return r.db.Save(file).Error
}

func (r *fileRepository) DeleteByUserID(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&entity.File{}).Error
}

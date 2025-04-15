package repository

import (
	"user_crud/internal/domain/entity"
	"user_crud/internal/domain/repository/interfaces"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) interfaces.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *entity.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindAll() ([]entity.User, error) {
	var users []entity.User
	err := r.db.Preload("Role").Find(&users).Error
	return users, err
}

func (r *userRepository) FindByID(id uint) (entity.User, error) {
	var user entity.User
	err := r.db.Preload("Role").First(&user, id).Error
	return user, err
}

func (r *userRepository) FindByEmail(email string) (entity.User, error) {
	var user entity.User
	err := r.db.Preload("Role").Where("email = ?", email).First(&user).Error
	return user, err
}

func (r *userRepository) Update(user *entity.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&entity.User{}, id).Error
}

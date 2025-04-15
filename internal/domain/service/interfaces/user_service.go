package interfaces

import (
	"user_crud/internal/dto"

	"github.com/gofiber/fiber/v2"
)

type UserService interface {
	CreateUser(c *fiber.Ctx, currentUserRole string) (dto.UserResponse, error, int)
	GetAllUsers(c *fiber.Ctx) ([]dto.UserResponse, error, int)
	GetUser(id uint, c *fiber.Ctx) (dto.UserResponse, error, int)
	UpdateUser(c *fiber.Ctx, id uint, currentUserID uint, currentUserRole string) (dto.UserResponse, error, int)
	DeleteUser(id uint, currentUserID uint, currentUserRole string) (error, int)
}

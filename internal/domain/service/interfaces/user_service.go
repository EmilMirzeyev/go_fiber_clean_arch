package interfaces

import (
	"user_crud/internal/dto"

	"github.com/gofiber/fiber/v2"
)

type UserService interface {
	CreateUser(c *fiber.Ctx) (dto.UserResponse, error)
	GetAllUsers(c *fiber.Ctx) ([]dto.UserResponse, error)
	GetUser(c *fiber.Ctx, id uint) (dto.UserResponse, error, int)
	UpdateUser(c *fiber.Ctx, id uint) (dto.UserResponse, error, int)
	DeleteUser(id uint) (error, int)
}

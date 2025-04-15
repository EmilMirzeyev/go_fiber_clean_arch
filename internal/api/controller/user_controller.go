package controller

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"user_crud/internal/domain/service/interfaces"
)

type UserController struct {
	userService interfaces.UserService
}

func NewUserController(userService interfaces.UserService) *UserController {
	return &UserController{
		userService,
	}
}

func (uc *UserController) CreateUser(c *fiber.Ctx) error {
	// Get current user role from context (set by auth middleware)
	role := c.Locals("role").(string)

	response, err, status := uc.userService.CreateUser(c, role)
	if err != nil {
		return fiber.NewError(status, err.Error())
	}

	return c.Status(status).JSON(response)
}

func (uc *UserController) GetAllUsers(c *fiber.Ctx) error {
	users, err, status := uc.userService.GetAllUsers(c)
	if err != nil {
		return fiber.NewError(status, err.Error())
	}

	return c.Status(status).JSON(users)
}

func (uc *UserController) GetUser(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid user ID")
	}

	user, err, status := uc.userService.GetUser(uint(id), c)
	if err != nil {
		return fiber.NewError(status, err.Error())
	}

	return c.Status(status).JSON(user)
}

func (uc *UserController) UpdateUser(c *fiber.Ctx) error {
	// Get target user ID
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid user ID")
	}

	// Get current user info from context (set by auth middleware)
	currentUserID := c.Locals("user_id").(uint)
	currentUserRole := c.Locals("role").(string)

	user, err, status := uc.userService.UpdateUser(c, uint(id), currentUserID, currentUserRole)
	if err != nil {
		return fiber.NewError(status, err.Error())
	}

	return c.Status(status).JSON(user)
}

func (uc *UserController) DeleteUser(c *fiber.Ctx) error {
	// Get target user ID
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid user ID")
	}

	// Get current user info from context (set by auth middleware)
	currentUserID := c.Locals("user_id").(uint)
	currentUserRole := c.Locals("role").(string)

	err, status := uc.userService.DeleteUser(uint(id), currentUserID, currentUserRole)
	if err != nil {
		return fiber.NewError(status, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

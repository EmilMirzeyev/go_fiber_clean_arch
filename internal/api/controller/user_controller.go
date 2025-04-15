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
	response, err := uc.userService.CreateUser(c)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

func (uc *UserController) GetAllUsers(c *fiber.Ctx) error {
	users, err := uc.userService.GetAllUsers(c)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(users)
}

func (uc *UserController) GetUser(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid user ID")
	}

	user, err, status := uc.userService.GetUser(c, uint(id))
	if err != nil {
		return fiber.NewError(status, err.Error())
	}

	return c.JSON(user)
}

func (uc *UserController) UpdateUser(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid user ID")
	}

	user, err, status := uc.userService.UpdateUser(c, uint(id))
	if err != nil {
		return fiber.NewError(status, err.Error())
	}

	return c.JSON(user)
}

func (uc *UserController) DeleteUser(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid user ID")
	}

	err, status := uc.userService.DeleteUser(uint(id))
	if err != nil {
		return fiber.NewError(status, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

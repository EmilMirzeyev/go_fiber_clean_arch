package controller

import (
	"github.com/gofiber/fiber/v2"

	"user_crud/internal/domain/service/interfaces"
	"user_crud/internal/dto"
)

type AuthController struct {
	authService interfaces.AuthService
}

func NewAuthController(authService interfaces.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

func (ac *AuthController) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	response, err, status := ac.authService.Register(req)
	if err != nil {
		return fiber.NewError(status, err.Error())
	}

	return c.Status(status).JSON(response)
}

func (ac *AuthController) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	response, err, status := ac.authService.Login(req)
	if err != nil {
		return fiber.NewError(status, err.Error())
	}

	return c.Status(status).JSON(response)
}

func (ac *AuthController) RefreshToken(c *fiber.Ctx) error {
	var req struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	response, err, status := ac.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		return fiber.NewError(status, err.Error())
	}

	return c.Status(status).JSON(response)
}

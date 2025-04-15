package routes

import (
	"github.com/gofiber/fiber/v2"

	"user_crud/internal/api/controller"
	"user_crud/internal/api/middleware"
)

func SetupRoutes(app *fiber.App, userController *controller.UserController, authController *controller.AuthController) {
	// Serve static files from public directory
	app.Static("/images", "./public/images")

	// API routes
	api := app.Group("/api")

	// Auth routes (public)
	auth := api.Group("/auth")
	auth.Post("/register", authController.Register)
	auth.Post("/login", authController.Login)
	auth.Post("/refresh", authController.RefreshToken)

	// User routes (protected)
	users := api.Group("/users", middleware.Protected())
	users.Post("/", middleware.RoleRequired("admin"), userController.CreateUser)
	users.Get("/", userController.GetAllUsers)
	users.Get("/:id", userController.GetUser)
	users.Put("/:id", userController.UpdateUser)
	users.Delete("/:id", middleware.RoleRequired("admin"), userController.DeleteUser)
}

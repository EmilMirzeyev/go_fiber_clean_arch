package routes

import (
	"github.com/gofiber/fiber/v2"

	"user_crud/internal/api/controller"
)

func SetupRoutes(app *fiber.App, userController *controller.UserController) {
	// Serve static files from public directory
	app.Static("/images", "./public/images")

	// API routes
	api := app.Group("/api")

	// User routes
	users := api.Group("/users")
	users.Post("/", userController.CreateUser)
	users.Get("/", userController.GetAllUsers)
	users.Get("/:id", userController.GetUser)
	users.Put("/:id", userController.UpdateUser)
	users.Delete("/:id", userController.DeleteUser)
}

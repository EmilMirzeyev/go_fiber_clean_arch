package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"user_crud/internal/api/controller"
	"user_crud/internal/api/routes"
	"user_crud/internal/config"
	"user_crud/internal/domain/repository"
	"user_crud/internal/domain/service"
	"user_crud/pkg/storage"
)

func main() {
	// Load configuration
	cfg := config.NewConfig()

	// Setup database connection
	db := storage.NewDatabaseConnection(cfg.DatabaseDSN)

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	fileRepo := repository.NewFileRepository(db)
	roleRepo := repository.NewRoleRepository(db)

	// Initialize services
	userService := service.NewUserService(userRepo, fileRepo)
	authService := service.NewAuthService(userRepo, roleRepo)

	// Initialize controllers
	userController := controller.NewUserController(userService)
	authController := controller.NewAuthController(authService)

	// Setup Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// Handle fiber-specific errors
			if fiberErr, ok := err.(*fiber.Error); ok {
				return c.Status(fiberErr.Code).JSON(fiber.Map{
					"error": fiberErr.Message,
				})
			}

			// Handle generic errors
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal Server Error",
			})
		},
	})

	// Add middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New())

	// Setup routes
	routes.SetupRoutes(app, userController, authController)

	// Start server
	serverAddr := fmt.Sprintf("%s:%d", cfg.ServerHost, cfg.ServerPort)
	log.Printf("Server is starting on %s\n", serverAddr)
	if err := app.Listen(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

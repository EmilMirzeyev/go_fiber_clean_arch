package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"user_crud/internal/util"
)

// Protected middleware to verify JWT access tokens
func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "Authorization header required")
		}

		// Check if the header has the Bearer prefix
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid authorization format")
		}

		// Extract the token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Verify the token
		claims, err := util.VerifyAccessToken(tokenString)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid or expired token")
		}

		// Set user information in context for later use
		c.Locals("user_id", claims.UserID)
		c.Locals("email", claims.Email)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}

// RoleRequired middleware to check if user has required role
func RoleRequired(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check if authenticated
		userRole := c.Locals("role")
		if userRole == nil {
			return fiber.NewError(fiber.StatusUnauthorized, "Authentication required")
		}

		// Check if user has one of the required roles
		roleStr := userRole.(string)
		for _, role := range roles {
			if role == roleStr {
				return c.Next()
			}
		}

		return fiber.NewError(fiber.StatusForbidden, "Insufficient permissions")
	}
}

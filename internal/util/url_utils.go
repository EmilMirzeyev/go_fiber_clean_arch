package util

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// BuildImageURL generates a full URL for an image using the request's host and protocol
func BuildImageURL(c *fiber.Ctx, imageName string) string {
	protocol := "http"
	if c.Protocol() == "https" {
		protocol = "https"
	}

	return fmt.Sprintf("%s://%s/images/%s", protocol, c.Hostname(), imageName)
}

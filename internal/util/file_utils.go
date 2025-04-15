package util

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
)

const (
	ImageDir = "public/images"
)

// SaveUploadedFile saves an uploaded file to the server's file system
// Changed the parameter type to multipart.FileHeader instead of fiber.FormFile
func SaveUploadedFile(c *fiber.Ctx, file *multipart.FileHeader) (string, error) {
	// Create unique filename
	imageName := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)

	// Ensure directory exists
	if err := os.MkdirAll(ImageDir, os.ModePerm); err != nil {
		return "", err
	}

	// Save the file
	savePath := filepath.Join(ImageDir, imageName)
	if err := c.SaveFile(file, savePath); err != nil {
		return "", err
	}

	return imageName, nil
}

// DeleteFile removes a file from the file system
func DeleteFile(filename string) error {
	if filename == "" {
		return nil
	}

	filePath := filepath.Join(ImageDir, filename)
	return os.Remove(filePath)
}

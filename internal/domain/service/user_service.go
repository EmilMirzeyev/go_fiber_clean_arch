package service

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"user_crud/internal/domain/entity"
	"user_crud/internal/domain/repository/interfaces"
	serviceInterfaces "user_crud/internal/domain/service/interfaces"
	"user_crud/internal/dto"
	"user_crud/internal/util"
)

type userService struct {
	userRepo interfaces.UserRepository
	fileRepo interfaces.FileRepository
}

func NewUserService(
	userRepo interfaces.UserRepository,
	fileRepo interfaces.FileRepository,
) serviceInterfaces.UserService {
	return &userService{
		userRepo: userRepo,
		fileRepo: fileRepo,
	}
}

func (s *userService) CreateUser(c *fiber.Ctx, currentUserRole string) (dto.UserResponse, error, int) {
	// Only admins can create users
	if currentUserRole != "admin" {
		return dto.UserResponse{}, errors.New("permission denied"), fiber.StatusForbidden
	}

	// Check required fields
	name := c.FormValue("name")
	if name == "" {
		return dto.UserResponse{}, errors.New("name is required"), fiber.StatusBadRequest
	}

	birthdate := c.FormValue("birthdate")
	if birthdate == "" {
		return dto.UserResponse{}, errors.New("birthdate is required"), fiber.StatusBadRequest
	}

	image, err := c.FormFile("image")
	if err != nil {
		return dto.UserResponse{}, errors.New("image is required"), fiber.StatusBadRequest
	}

	// Parse birthdate
	birthTime, err := util.ParseBirthdate(birthdate)
	if err != nil {
		return dto.UserResponse{}, errors.New("invalid birthdate format. Please use DD.MM.YYYY"), fiber.StatusBadRequest
	}

	// Calculate age
	age := util.CalculateAge(birthTime)

	// Save image file
	imageName, err := util.SaveUploadedFile(c, image)
	if err != nil {
		return dto.UserResponse{}, fmt.Errorf("failed to save image: %w", err), fiber.StatusInternalServerError
	}

	// Create user entity
	user := entity.User{
		Name:      name,
		Age:       age,
		ImageName: imageName,
	}

	// Create user in database
	if err := s.userRepo.Create(&user); err != nil {
		// Clean up the image file if user creation fails
		_ = util.DeleteFile(imageName)
		return dto.UserResponse{}, fmt.Errorf("failed to create user: %w", err), fiber.StatusInternalServerError
	}

	// Create file record
	file := entity.File{
		FileName: imageName,
		UserID:   user.ID,
	}

	if err := s.fileRepo.Create(&file); err != nil {
		// We should implement a proper rollback here in production code
		return dto.UserResponse{}, fmt.Errorf("failed to create file record: %w", err), fiber.StatusInternalServerError
	}

	// Build response
	return dto.UserResponse{
		ID:       user.ID,
		Name:     user.Name,
		Age:      user.Age,
		Email:    user.Email,
		Role:     user.Role.Name,
		ImageUrl: util.BuildImageURL(c, user.ImageName),
	}, nil, fiber.StatusCreated
}

func (s *userService) GetAllUsers(c *fiber.Ctx) ([]dto.UserResponse, error, int) {
	users, err := s.userRepo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve users: %w", err), fiber.StatusInternalServerError
	}

	var response []dto.UserResponse
	for _, user := range users {
		response = append(response, dto.UserResponse{
			ID:       user.ID,
			Name:     user.Name,
			Age:      user.Age,
			Email:    user.Email,
			Role:     user.Role.Name,
			ImageUrl: util.BuildImageURL(c, user.ImageName),
		})
	}

	return response, nil, fiber.StatusOK
}

func (s *userService) GetUser(id uint, c *fiber.Ctx) (dto.UserResponse, error, int) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.UserResponse{}, errors.New("user not found"), fiber.StatusNotFound
		}
		return dto.UserResponse{}, fmt.Errorf("failed to retrieve user: %w", err), fiber.StatusInternalServerError
	}

	return dto.UserResponse{
		ID:       user.ID,
		Name:     user.Name,
		Age:      user.Age,
		Email:    user.Email,
		Role:     user.Role.Name,
		ImageUrl: util.BuildImageURL(c, user.ImageName),
	}, nil, fiber.StatusOK
}

func (s *userService) UpdateUser(c *fiber.Ctx, id uint, currentUserID uint, currentUserRole string) (dto.UserResponse, error, int) {
	// Check if user has permission to update this record
	// Admin can update any record, users can only update their own
	if currentUserRole != "admin" && currentUserID != id {
		return dto.UserResponse{}, errors.New("permission denied"), fiber.StatusForbidden
	}

	// Find existing user
	existingUser, err := s.userRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.UserResponse{}, errors.New("user not found"), fiber.StatusNotFound
		}
		return dto.UserResponse{}, fmt.Errorf("failed to find user: %w", err), fiber.StatusInternalServerError
	}

	// Check required fields
	name := c.FormValue("name")
	if name == "" {
		return dto.UserResponse{}, errors.New("name is required"), fiber.StatusBadRequest
	}

	birthdate := c.FormValue("birthdate")
	if birthdate == "" {
		return dto.UserResponse{}, errors.New("birthdate is required"), fiber.StatusBadRequest
	}

	// Parse birthdate and update age
	birthTime, err := util.ParseBirthdate(birthdate)
	if err != nil {
		return dto.UserResponse{}, errors.New("invalid birthdate format. Please use DD.MM.YYYY"), fiber.StatusBadRequest
	}

	// Update user fields
	existingUser.Name = name
	existingUser.Age = util.CalculateAge(birthTime)

	// Check if there's a new image
	oldImageName := existingUser.ImageName
	image, err := c.FormFile("image")

	if err == nil {
		// New image was uploaded
		imageName, err := util.SaveUploadedFile(c, image)
		if err != nil {
			return dto.UserResponse{}, fmt.Errorf("failed to save new image: %w", err), fiber.StatusInternalServerError
		}

		// Update user with new image
		existingUser.ImageName = imageName

		// Update or create file record
		file, err := s.fileRepo.FindByUserID(existingUser.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Create new file record
				newFile := entity.File{
					FileName: imageName,
					UserID:   existingUser.ID,
				}
				if err := s.fileRepo.Create(&newFile); err != nil {
					return dto.UserResponse{}, fmt.Errorf("failed to create file record: %w", err), fiber.StatusInternalServerError
				}
			} else {
				return dto.UserResponse{}, fmt.Errorf("failed to retrieve file record: %w", err), fiber.StatusInternalServerError
			}
		} else {
			// Update existing file
			file.FileName = imageName
			if err := s.fileRepo.Update(&file); err != nil {
				return dto.UserResponse{}, fmt.Errorf("failed to update file record: %w", err), fiber.StatusInternalServerError
			}
		}

		// Delete old image
		_ = util.DeleteFile(oldImageName)
	}

	// Save updated user
	if err := s.userRepo.Update(&existingUser); err != nil {
		return dto.UserResponse{}, fmt.Errorf("failed to update user: %w", err), fiber.StatusInternalServerError
	}

	return dto.UserResponse{
		ID:       existingUser.ID,
		Name:     existingUser.Name,
		Age:      existingUser.Age,
		Email:    existingUser.Email,
		Role:     existingUser.Role.Name,
		ImageUrl: util.BuildImageURL(c, existingUser.ImageName),
	}, nil, fiber.StatusOK
}

func (s *userService) DeleteUser(id uint, currentUserID uint, currentUserRole string) (error, int) {
	// Only admins can delete users
	if currentUserRole != "admin" {
		return errors.New("permission denied"), fiber.StatusForbidden
	}

	// Find the user to get image filename
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found"), fiber.StatusNotFound
		}
		return fmt.Errorf("failed to find user: %w", err), fiber.StatusInternalServerError
	}

	// Store image name for later deletion
	imageName := user.ImageName

	// Delete file records first (respect foreign key constraints)
	if err := s.fileRepo.DeleteByUserID(id); err != nil {
		return fmt.Errorf("failed to delete file records: %w", err), fiber.StatusInternalServerError
	}

	// Delete the user
	if err := s.userRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err), fiber.StatusInternalServerError
	}

	// Delete the image file
	_ = util.DeleteFile(imageName)

	return nil, fiber.StatusNoContent
}

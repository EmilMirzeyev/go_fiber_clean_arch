package service

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"user_crud/internal/domain/entity"
	"user_crud/internal/domain/repository/interfaces"
	serviceInterfaces "user_crud/internal/domain/service/interfaces"
	"user_crud/internal/dto"
	"user_crud/internal/util"
)

type authService struct {
	userRepo interfaces.UserRepository
	roleRepo interfaces.RoleRepository
}

func NewAuthService(
	userRepo interfaces.UserRepository,
	roleRepo interfaces.RoleRepository,
) serviceInterfaces.AuthService {
	return &authService{
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

func (s *authService) Register(req dto.RegisterRequest) (dto.TokenResponse, error, int) {
	// Check if email already exists
	_, err := s.userRepo.FindByEmail(req.Email)
	if err == nil {
		return dto.TokenResponse{}, errors.New("email already exists"), fiber.StatusConflict
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return dto.TokenResponse{}, errors.New("failed to check email existence"), fiber.StatusInternalServerError
	}

	// Hash password
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return dto.TokenResponse{}, errors.New("failed to hash password"), fiber.StatusInternalServerError
	}

	// Get default user role
	role, err := s.roleRepo.FindByName("user")
	if err != nil {
		return dto.TokenResponse{}, errors.New("default role not found"), fiber.StatusInternalServerError
	}

	// Create user
	user := entity.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
		RoleID:   role.ID,
		Age:      0, // Set default age
	}

	if err := s.userRepo.Create(&user); err != nil {
		return dto.TokenResponse{}, errors.New("failed to create user"), fiber.StatusInternalServerError
	}

	// Generate tokens
	accessToken, err := util.GenerateAccessToken(user.ID, user.Email, role.Name)
	if err != nil {
		return dto.TokenResponse{}, errors.New("failed to generate access token"), fiber.StatusInternalServerError
	}

	refreshToken, err := util.GenerateRefreshToken(user.ID)
	if err != nil {
		return dto.TokenResponse{}, errors.New("failed to generate refresh token"), fiber.StatusInternalServerError
	}

	return dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(util.AccessTokenExpiry / time.Second),
	}, nil, fiber.StatusCreated
}

func (s *authService) Login(req dto.LoginRequest) (dto.TokenResponse, error, int) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.TokenResponse{}, errors.New("invalid email or password"), fiber.StatusUnauthorized
		}
		return dto.TokenResponse{}, errors.New("failed to retrieve user"), fiber.StatusInternalServerError
	}

	// Verify password
	if !util.CheckPassword(req.Password, user.Password) {
		return dto.TokenResponse{}, errors.New("invalid email or password"), fiber.StatusUnauthorized
	}

	// Get role name
	role, err := s.roleRepo.FindByID(user.RoleID)
	if err != nil {
		return dto.TokenResponse{}, errors.New("failed to retrieve role"), fiber.StatusInternalServerError
	}

	// Generate tokens
	accessToken, err := util.GenerateAccessToken(user.ID, user.Email, role.Name)
	if err != nil {
		return dto.TokenResponse{}, errors.New("failed to generate access token"), fiber.StatusInternalServerError
	}

	refreshToken, err := util.GenerateRefreshToken(user.ID)
	if err != nil {
		return dto.TokenResponse{}, errors.New("failed to generate refresh token"), fiber.StatusInternalServerError
	}

	return dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(util.AccessTokenExpiry / time.Second),
	}, nil, fiber.StatusOK
}

func (s *authService) RefreshToken(refreshToken string) (dto.TokenResponse, error, int) {
	// Verify refresh token
	userID, err := util.VerifyRefreshToken(refreshToken)
	if err != nil {
		return dto.TokenResponse{}, errors.New("invalid refresh token"), fiber.StatusUnauthorized
	}

	// Get user
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return dto.TokenResponse{}, errors.New("user not found"), fiber.StatusUnauthorized
	}

	// Get role
	role, err := s.roleRepo.FindByID(user.RoleID)
	if err != nil {
		return dto.TokenResponse{}, errors.New("failed to retrieve role"), fiber.StatusInternalServerError
	}

	// Generate new tokens
	accessToken, err := util.GenerateAccessToken(user.ID, user.Email, role.Name)
	if err != nil {
		return dto.TokenResponse{}, errors.New("failed to generate access token"), fiber.StatusInternalServerError
	}

	newRefreshToken, err := util.GenerateRefreshToken(user.ID)
	if err != nil {
		return dto.TokenResponse{}, errors.New("failed to generate refresh token"), fiber.StatusInternalServerError
	}

	return dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(util.AccessTokenExpiry / time.Second),
	}, nil, fiber.StatusOK
}

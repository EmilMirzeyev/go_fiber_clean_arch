package interfaces

import (
	"user_crud/internal/dto"
)

type AuthService interface {
	Register(req dto.RegisterRequest) (dto.TokenResponse, error, int)
	Login(req dto.LoginRequest) (dto.TokenResponse, error, int)
	RefreshToken(refreshToken string) (dto.TokenResponse, error, int)
}

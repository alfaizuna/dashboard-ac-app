package request

import "dashboard-ac-backend/internal/domain"

type RegisterRequest struct {
	Name     string      `json:"name" validate:"required,min=2,max=100"`
	Email    string      `json:"email" validate:"required,email"`
	Password string      `json:"password" validate:"required,min=6"`
	Role     domain.Role `json:"role" validate:"omitempty,oneof=admin technician customer"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
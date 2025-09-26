package handler

import (
	"dashboard-ac-backend/internal/api/request"
	"dashboard-ac-backend/internal/api/response"
	"dashboard-ac-backend/internal/domain"
	"dashboard-ac-backend/internal/service"
	"dashboard-ac-backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req request.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", nil)
	}

	// Validate request
	if errors := utils.ValidateStruct(&req); len(errors) > 0 {
		return response.BadRequest(c, "Validation failed", errors)
	}

	// Register user
	user, err := h.authService.Register(&req)
	if err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}

	return response.Created(c, "User registered successfully", map[string]interface{}{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"role":  user.Role,
	})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req request.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", nil)
	}

	// Validate request
	if errors := utils.ValidateStruct(&req); len(errors) > 0 {
		return response.BadRequest(c, "Validation failed", errors)
	}

	// Login user
	tokenPair, user, err := h.authService.Login(&req)
	if err != nil {
		return response.Unauthorized(c, err.Error())
	}

	return response.Success(c, "Login successful", map[string]interface{}{
		"user": map[string]interface{}{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
		"tokens": tokenPair,
	})
}

func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req request.RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", nil)
	}

	// Validate request
	if errors := utils.ValidateStruct(&req); len(errors) > 0 {
		return response.BadRequest(c, "Validation failed", errors)
	}

	// Refresh token
	tokenPair, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		return response.Unauthorized(c, err.Error())
	}

	return response.Success(c, "Token refreshed successfully", map[string]interface{}{
		"tokens": tokenPair,
	})
}

func (h *AuthHandler) Me(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	userEmail := c.Locals("user_email").(string)
	userRole := c.Locals("user_role").(domain.Role)

	return response.Success(c, "User profile retrieved successfully", map[string]interface{}{
		"id":    userID,
		"email": userEmail,
		"role":  string(userRole),
	})
}
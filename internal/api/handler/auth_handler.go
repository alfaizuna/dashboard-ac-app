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
    // Register godoc
    // @Summary Register user baru
    // @Description Mendaftarkan pengguna baru (role: admin/technician/customer)
    // @Tags Auth
    // @Accept json
    // @Produce json
    // @Param request body request.RegisterRequest true "Register Request"
    // @Success 201 {object} response.BaseResponse
    // @Failure 400 {object} response.BaseResponse
    // @Router /auth/register [post]
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

	// Prepare response data
	responseData := map[string]interface{}{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"role":  user.Role,
	}

	// Add success message based on role
	var message string
	if user.Role == domain.RoleCustomer {
		message = "Customer registered successfully"
	} else {
		message = "User registered successfully"
	}

	return response.Created(c, message, responseData)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
    // Login godoc
    // @Summary Login
    // @Description Autentikasi pengguna dengan email dan password
    // @Tags Auth
    // @Accept json
    // @Produce json
    // @Param request body request.LoginRequest true "Login Request"
    // @Success 200 {object} response.BaseResponse
    // @Failure 401 {object} response.BaseResponse
    // @Failure 400 {object} response.BaseResponse
    // @Router /auth/login [post]
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
    // RefreshToken godoc
    // @Summary Refresh token
    // @Description Menghasilkan pasangan token baru dari refresh token
    // @Tags Auth
    // @Accept json
    // @Produce json
    // @Param request body request.RefreshTokenRequest true "Refresh Token Request"
    // @Success 200 {object} response.BaseResponse
    // @Failure 401 {object} response.BaseResponse
    // @Failure 400 {object} response.BaseResponse
    // @Router /auth/refresh [post]
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
    // Me godoc
    // @Summary Profil pengguna saat ini
    // @Description Mendapatkan profil user dari token JWT
    // @Tags Auth
    // @Security BearerAuth
    // @Produce json
    // @Success 200 {object} response.BaseResponse
    // @Failure 401 {object} response.BaseResponse
    // @Router /me [get]
    userID := c.Locals("user_id").(uint)
    userEmail := c.Locals("user_email").(string)
    userRole := c.Locals("user_role").(domain.Role)

    return response.Success(c, "User profile retrieved successfully", map[string]interface{}{
        "id":    userID,
        "email": userEmail,
        "role":  string(userRole),
    })
}
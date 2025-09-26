package handler

import (
	"strconv"

	"dashboard-ac-backend/internal/api/request"
	"dashboard-ac-backend/internal/api/response"
	"dashboard-ac-backend/internal/domain"
	"dashboard-ac-backend/internal/service"
	"dashboard-ac-backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) Create(c *fiber.Ctx) error {
	var req request.UserCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", nil)
	}

	// Validate request
	if errors := utils.ValidateStruct(&req); len(errors) > 0 {
		return response.BadRequest(c, "Validation failed", errors)
	}

	// Create user
	user, err := h.userService.Create(&req)
	if err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}

	return response.Created(c, "User created successfully", map[string]interface{}{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"role":  user.Role,
	})
}

func (h *UserHandler) GetByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID", nil)
	}

	user, err := h.userService.GetByID(uint(id))
	if err != nil {
		return response.NotFound(c, err.Error())
	}

	return response.Success(c, "User retrieved successfully", map[string]interface{}{
		"id":         user.ID,
		"name":       user.Name,
		"email":      user.Email,
		"role":       user.Role,
		"is_active":  user.IsActive,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	})
}

func (h *UserHandler) Update(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID", nil)
	}

	var req request.UserUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", nil)
	}

	// Validate request
	if errors := utils.ValidateStruct(&req); len(errors) > 0 {
		return response.BadRequest(c, "Validation failed", errors)
	}

	// Update user
	user, err := h.userService.Update(uint(id), &req)
	if err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}

	return response.Success(c, "User updated successfully", map[string]interface{}{
		"id":         user.ID,
		"name":       user.Name,
		"email":      user.Email,
		"role":       user.Role,
		"is_active":  user.IsActive,
		"updated_at": user.UpdatedAt,
	})
}

func (h *UserHandler) Delete(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID", nil)
	}

	if err := h.userService.Delete(uint(id)); err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}

	return response.Success(c, "User deleted successfully", nil)
}

func (h *UserHandler) List(c *fiber.Ctx) error {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	pagination := &request.PaginationRequest{
		Page:  page,
		Limit: limit,
	}

	// Get users
	users, total, err := h.userService.List(pagination)
	if err != nil {
		return response.InternalServerError(c, "Failed to retrieve users")
	}

	// Format response data
	var userData []map[string]interface{}
	for _, user := range users {
		userData = append(userData, map[string]interface{}{
			"id":         user.ID,
			"name":       user.Name,
			"email":      user.Email,
			"role":       user.Role,
			"is_active":  user.IsActive,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		})
	}

	// Calculate pagination metadata
	paginationMeta := response.CalculatePagination(pagination.Page, pagination.GetLimit(), total)

	return response.Paginated(c, "Users retrieved successfully", userData, paginationMeta)
}

func (h *UserHandler) GetByRole(c *fiber.Ctx) error {
	roleParam := c.Params("role")
	role := domain.Role(roleParam)

	// Validate role
	if role != domain.RoleAdmin && role != domain.RoleTechnician && role != domain.RoleCustomer {
		return response.BadRequest(c, "Invalid role", nil)
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	pagination := &request.PaginationRequest{
		Page:  page,
		Limit: limit,
	}

	// Get users by role
	users, total, err := h.userService.GetByRole(role, pagination)
	if err != nil {
		return response.InternalServerError(c, "Failed to retrieve users")
	}

	// Format response data
	var userData []map[string]interface{}
	for _, user := range users {
		userData = append(userData, map[string]interface{}{
			"id":         user.ID,
			"name":       user.Name,
			"email":      user.Email,
			"role":       user.Role,
			"is_active":  user.IsActive,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		})
	}

	// Calculate pagination metadata
	paginationMeta := response.CalculatePagination(pagination.Page, pagination.GetLimit(), total)

	return response.Paginated(c, "Users retrieved successfully", userData, paginationMeta)
}
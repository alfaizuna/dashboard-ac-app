package request

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type PaginationRequest struct {
	Page   int    `json:"page" query:"page" validate:"min=1"`
	Limit  int    `json:"limit" query:"limit" validate:"min=1,max=100"`
	Sort   string `json:"sort" query:"sort"`
	Filter string `json:"filter" query:"filter"`
	Search string `json:"search" query:"search"`
}

// GetPaginationFromQuery extracts pagination parameters from query string
func GetPaginationFromQuery(c *fiber.Ctx) *PaginationRequest {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	return &PaginationRequest{
		Page:   page,
		Limit:  limit,
		Sort:   c.Query("sort", "created_at desc"),
		Filter: c.Query("filter", ""),
		Search: c.Query("search", ""),
	}
}

func (p *PaginationRequest) GetOffset() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	return (p.Page - 1) * p.GetLimit()
}

func (p *PaginationRequest) GetLimit() int {
	if p.Limit <= 0 {
		p.Limit = 10
	}
	if p.Limit > 100 {
		p.Limit = 100
	}
	return p.Limit
}

type UserCreateRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Role     string `json:"role" validate:"required,oneof=admin technician customer"`
}

type UserUpdateRequest struct {
	Name     *string `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Email    *string `json:"email,omitempty" validate:"omitempty,email"`
	Role     *string `json:"role,omitempty" validate:"omitempty,oneof=admin technician customer"`
	IsActive *bool   `json:"is_active,omitempty"`
}
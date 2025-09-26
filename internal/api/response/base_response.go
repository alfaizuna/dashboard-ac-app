package response

import (
	"github.com/gofiber/fiber/v2"
)

type BaseResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type PaginatedResponse struct {
	Status     string         `json:"status"`
	Message    string         `json:"message"`
	Data       interface{}    `json:"data"`
	Pagination PaginationMeta `json:"pagination"`
}

type ErrorDetail struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
}

func Success(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(BaseResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

func Created(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(BaseResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

func BadRequest(c *fiber.Ctx, message string, errors interface{}) error {
	return c.Status(fiber.StatusBadRequest).JSON(BaseResponse{
		Status:  "error",
		Message: message,
		Error:   errors,
	})
}

func Unauthorized(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusUnauthorized).JSON(BaseResponse{
		Status:  "error",
		Message: message,
	})
}

func Forbidden(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusForbidden).JSON(BaseResponse{
		Status:  "error",
		Message: message,
	})
}

func NotFound(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusNotFound).JSON(BaseResponse{
		Status:  "error",
		Message: message,
	})
}

func InternalServerError(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusInternalServerError).JSON(BaseResponse{
		Status:  "error",
		Message: message,
	})
}

func Paginated(c *fiber.Ctx, message string, data interface{}, pagination PaginationMeta) error {
	return c.Status(fiber.StatusOK).JSON(PaginatedResponse{
		Status:     "success",
		Message:    message,
		Data:       data,
		Pagination: pagination,
	})
}

func CalculatePagination(page, limit int, total int64) PaginationMeta {
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	return PaginationMeta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}
}
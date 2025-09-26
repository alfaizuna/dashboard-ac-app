package handler

import (
	"time"

	"dashboard-ac-backend/internal/api/response"

	"github.com/gofiber/fiber/v2"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Check(c *fiber.Ctx) error {
	return response.Success(c, "Server is healthy", map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now().UTC(),
		"service":   "dashboard-ac-backend",
		"version":   "1.0.0",
	})
}
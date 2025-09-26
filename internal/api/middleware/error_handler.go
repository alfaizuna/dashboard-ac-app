package middleware

import (
	"dashboard-ac-backend/internal/api/response"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	// Default error code
	code := fiber.StatusInternalServerError

	// Check if it's a Fiber error
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	// Log the error
	log.Error().
		Err(err).
		Str("method", c.Method()).
		Str("path", c.Path()).
		Str("ip", c.IP()).
		Int("status", code).
		Msg("HTTP Error")

	// Handle different error codes
	switch code {
	case fiber.StatusBadRequest:
		return response.BadRequest(c, err.Error(), nil)
	case fiber.StatusUnauthorized:
		return response.Unauthorized(c, err.Error())
	case fiber.StatusForbidden:
		return response.Forbidden(c, err.Error())
	case fiber.StatusNotFound:
		return response.NotFound(c, err.Error())
	default:
		return response.InternalServerError(c, "Internal server error")
	}
}
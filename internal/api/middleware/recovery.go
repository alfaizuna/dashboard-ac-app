package middleware

import (
	"dashboard-ac-backend/internal/api/response"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog/log"
)

func Recovery() fiber.Handler {
	return recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
			// Log the panic with stack trace
			log.Error().
				Interface("panic", e).
				Str("method", c.Method()).
				Str("path", c.Path()).
				Str("ip", c.IP()).
				Msg("Panic recovered")

			// Return internal server error response
			response.InternalServerError(c, "Internal server error")
		},
	})
}
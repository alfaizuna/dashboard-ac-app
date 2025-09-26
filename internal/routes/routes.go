package routes

import (
	"dashboard-ac-backend/internal/api/handler"
	"dashboard-ac-backend/internal/api/middleware"
	"dashboard-ac-backend/internal/service"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, authService service.AuthService, userService service.UserService, jwtSecret string) {
	// Initialize handlers
	healthHandler := handler.NewHealthHandler()
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)

	// Health check endpoint
	app.Get("/health", healthHandler.Check)

	// API v1 routes
	api := app.Group("/api/v1")

	// Auth routes (public)
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.RefreshToken)

	// Protected routes
	protected := api.Group("", middleware.JWTAuth(jwtSecret))
	
	// User profile
	protected.Get("/me", authHandler.Me)

	// User management routes (admin only)
	users := protected.Group("/users", middleware.RequireAdmin())
	users.Post("/", userHandler.Create)
	users.Get("/", userHandler.List)
	users.Get("/:id", userHandler.GetByID)
	users.Put("/:id", userHandler.Update)
	users.Delete("/:id", userHandler.Delete)
	users.Get("/role/:role", userHandler.GetByRole)

	// Technician routes (admin and technician)
	technicians := protected.Group("/technicians", middleware.RequireAdminOrTechnician())
	technicians.Get("/", func(c *fiber.Ctx) error {
		// This would be implemented later for technician-specific operations
		return c.JSON(fiber.Map{
			"message": "Technician routes - to be implemented",
		})
	})
}
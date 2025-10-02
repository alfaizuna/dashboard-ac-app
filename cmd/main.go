package main

import (
	"log"

	"dashboard-ac-backend/config"
	"dashboard-ac-backend/internal/api/middleware"
	"dashboard-ac-backend/internal/repository"
	"dashboard-ac-backend/internal/routes"
	"dashboard-ac-backend/internal/service"
	"dashboard-ac-backend/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Initialize logger
	logger.InitLogger(cfg.Environment)

	// Initialize database
	db, err := config.InitDatabase(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	technicianRepo := repository.NewTechnicianRepository(db)
	serviceRepo := repository.NewServiceRepository(db)
	scheduleRepo := repository.NewScheduleRepository(db)
	invoiceRepo := repository.NewInvoiceRepository(db)
	invoiceDetailRepo := repository.NewInvoiceDetailRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	userService := service.NewUserService(userRepo)
	customerService := service.NewCustomerService(customerRepo)
	technicianService := service.NewTechnicianService(technicianRepo)
	serviceService := service.NewServiceService(serviceRepo)
	scheduleService := service.NewScheduleService(scheduleRepo, customerRepo, technicianRepo, serviceRepo)
	invoiceService := service.NewInvoiceService(invoiceRepo, customerRepo, scheduleRepo)
	invoiceDetailService := service.NewInvoiceDetailService(invoiceDetailRepo, invoiceRepo, serviceRepo)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})

	// Setup routes with all services
	routes.SetupRoutes(
		app,
		authService,
		userService,
		customerService,
		technicianService,
		serviceService,
		scheduleService,
		invoiceService,
		invoiceDetailService,
		cfg.JWTSecret,
	)

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
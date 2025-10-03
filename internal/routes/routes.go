package routes

import (
    "dashboard-ac-backend/internal/api/handler"
    "dashboard-ac-backend/internal/api/middleware"
    "dashboard-ac-backend/internal/service"

    "github.com/gofiber/fiber/v2"
    swagger "github.com/gofiber/swagger"
)

func SetupRoutes(
	app *fiber.App,
	authService service.AuthService,
	userService service.UserService,
	customerService service.CustomerService,
	technicianService service.TechnicianService,
	serviceService service.ServiceService,
	scheduleService service.ScheduleService,
	invoiceService service.InvoiceService,
	invoiceDetailService service.InvoiceDetailService,
	jwtSecret string,
) {
    // Setup global middleware
    app.Use(middleware.Logger())
    app.Use(middleware.Recovery())
    app.Use(middleware.CORS())

    // Swagger UI route (accessible without auth)
    // Visit http://localhost:<port>/docs to view API documentation
    app.Get("/docs/*", swagger.New(swagger.Config{
        Title: "Dashboard AC Backend API",
        DeepLinking: true,
        // HideTop field removed as it's not supported in swagger.Config
    }))

    // Initialize handlers
    healthHandler := handler.NewHealthHandler()
    authHandler := handler.NewAuthHandler(authService)
    userHandler := handler.NewUserHandler(userService)
    customerHandler := handler.NewCustomerHandler(customerService)
	technicianHandler := handler.NewTechnicianHandler(technicianService)
	serviceHandler := handler.NewServiceHandler(serviceService)
	scheduleHandler := handler.NewScheduleHandler(scheduleService)
	invoiceHandler := handler.NewInvoiceHandler(invoiceService)
	invoiceDetailHandler := handler.NewInvoiceDetailHandler(invoiceDetailService)

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

	// Customer management routes (admin and technician)
	customers := protected.Group("/customers", middleware.RequireAdminOrTechnician())
	customers.Post("/", customerHandler.CreateCustomer)
	customers.Get("/", customerHandler.ListCustomers)
	customers.Get("/:id", customerHandler.GetCustomer)
	customers.Put("/:id", customerHandler.UpdateCustomer)
	customers.Delete("/:id", customerHandler.DeleteCustomer)
	customers.Get("/search", customerHandler.SearchCustomers)

	// Technician management routes (admin only)
	technicians := protected.Group("/technicians", middleware.RequireAdmin())
	technicians.Post("/", technicianHandler.CreateTechnician)
	technicians.Get("/", technicianHandler.ListTechnicians)
	technicians.Get("/:id", technicianHandler.GetTechnician)
	technicians.Put("/:id", technicianHandler.UpdateTechnician)
	technicians.Delete("/:id", technicianHandler.DeleteTechnician)
	technicians.Get("/search", technicianHandler.SearchTechnicians)

	// Service management routes (admin only)
	services := protected.Group("/services", middleware.RequireAdmin())
	services.Post("/", serviceHandler.CreateService)
	services.Get("/", serviceHandler.ListServices)
	services.Get("/:id", serviceHandler.GetService)
	services.Put("/:id", serviceHandler.UpdateService)
	services.Delete("/:id", serviceHandler.DeleteService)
	services.Get("/search", serviceHandler.SearchServices)

	// Schedule management routes (admin and technician)
	schedules := protected.Group("/schedules", middleware.RequireAdminOrTechnician())
	schedules.Post("/", scheduleHandler.CreateSchedule)
	schedules.Get("/", scheduleHandler.ListSchedules)
	schedules.Get("/:id", scheduleHandler.GetSchedule)
	schedules.Put("/:id", scheduleHandler.UpdateSchedule)
	schedules.Delete("/:id", scheduleHandler.DeleteSchedule)
	schedules.Get("/search", scheduleHandler.SearchSchedules)
	schedules.Get("/customer/:customer_id", scheduleHandler.GetSchedulesByCustomer)
	schedules.Get("/technician/:technician_id", scheduleHandler.GetSchedulesByTechnician)
	schedules.Get("/status/:status", scheduleHandler.GetSchedulesByStatus)

	// Invoice management routes (admin and technician)
	invoices := protected.Group("/invoices", middleware.RequireAdminOrTechnician())
	invoices.Post("/", invoiceHandler.CreateInvoice)
	invoices.Get("/", invoiceHandler.ListInvoices)
	invoices.Get("/:id", invoiceHandler.GetInvoice)
	invoices.Put("/:id", invoiceHandler.UpdateInvoice)
	invoices.Delete("/:id", invoiceHandler.DeleteInvoice)
	invoices.Get("/search", invoiceHandler.SearchInvoices)
	invoices.Get("/customer/:customer_id", invoiceHandler.GetInvoicesByCustomer)
	invoices.Get("/schedule/:schedule_id", invoiceHandler.GetInvoicesBySchedule)
	invoices.Get("/status/:status", invoiceHandler.GetInvoicesByStatus)

	// Invoice detail management routes (admin and technician)
	invoiceDetails := protected.Group("/invoice-details", middleware.RequireAdminOrTechnician())
	invoiceDetails.Post("/", invoiceDetailHandler.CreateInvoiceDetail)
	invoiceDetails.Get("/:id", invoiceDetailHandler.GetInvoiceDetail)
	invoiceDetails.Put("/:id", invoiceDetailHandler.UpdateInvoiceDetail)
	invoiceDetails.Delete("/:id", invoiceDetailHandler.DeleteInvoiceDetail)
	invoiceDetails.Get("/invoice/:invoice_id", invoiceDetailHandler.GetInvoiceDetailsByInvoice)
	invoiceDetails.Delete("/invoice/:invoice_id", invoiceDetailHandler.DeleteInvoiceDetailsByInvoice)
}
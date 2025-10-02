package main

import (
	"log"

	"dashboard-ac-backend/config"
	"dashboard-ac-backend/internal/domain"
	"dashboard-ac-backend/pkg/hash"

	"gorm.io/gorm"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Initialize database
	db, err := config.InitDatabase(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run migrations
	if err := runMigrations(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Seed database
	if err := seedDatabase(db); err != nil {
		log.Fatal("Failed to seed database:", err)
	}

	log.Println("Database migration and seeding completed successfully!")
}

func runMigrations(db *gorm.DB) error {
	log.Println("Running database migrations...")
	
	// Auto migrate all models
	return db.AutoMigrate(
		&domain.User{},
		&domain.Customer{},
		&domain.Technician{},
		&domain.Service{},
		&domain.Schedule{},
		&domain.Invoice{},
		&domain.InvoiceDetail{},
	)
}

func seedDatabase(db *gorm.DB) error {
	log.Println("Seeding database...")

	// Check if admin user already exists
	var adminUser domain.User
	result := db.Where("email = ?", "admin@example.com").First(&adminUser)
	if result.Error == nil {
		log.Println("Admin user already exists, skipping seed")
		return nil
	}

	// Hash default password
	hashedPassword, err := hash.HashPassword("admin123")
	if err != nil {
		return err
	}

	// Create default admin user
	admin := domain.User{
		Name:     "System Administrator",
		Email:    "admin@example.com",
		Password: hashedPassword,
		Role:     domain.RoleAdmin,
		IsActive: true,
	}

	if err := db.Create(&admin).Error; err != nil {
		return err
	}

	log.Printf("Default admin user created with email: %s and password: admin123", admin.Email)

	// Create default technician user
	techPassword, err := hash.HashPassword("tech123")
	if err != nil {
		return err
	}

	technician := domain.User{
		Name:     "Default Technician",
		Email:    "technician@example.com",
		Password: techPassword,
		Role:     domain.RoleTechnician,
		IsActive: true,
	}

	if err := db.Create(&technician).Error; err != nil {
		return err
	}

	log.Printf("Default technician user created with email: %s and password: tech123", technician.Email)

	// Create default customer user
	customerPassword, err := hash.HashPassword("customer123")
	if err != nil {
		return err
	}

	customer := domain.User{
		Name:     "Default Customer",
		Email:    "customer@example.com",
		Password: customerPassword,
		Role:     domain.RoleCustomer,
		IsActive: true,
	}

	if err := db.Create(&customer).Error; err != nil {
		return err
	}

	log.Printf("Default customer user created with email: %s and password: customer123", customer.Email)

	// Seed default services
	services := []domain.Service{
		{
			Name:     "Cuci AC",
			Price:    150000,
			Duration: 60, // 1 hour
		},
		{
			Name:     "Isi Freon",
			Price:    200000,
			Duration: 45, // 45 minutes
		},
		{
			Name:     "Bongkar Pasang AC",
			Price:    500000,
			Duration: 180, // 3 hours
		},
		{
			Name:     "Service Rutin AC",
			Price:    100000,
			Duration: 30, // 30 minutes
		},
	}

	for _, service := range services {
		var existingService domain.Service
		result := db.Where("name = ?", service.Name).First(&existingService)
		if result.Error != nil {
			if err := db.Create(&service).Error; err != nil {
				return err
			}
			log.Printf("Service created: %s - Rp %.0f", service.Name, service.Price)
		}
	}

	log.Println("Database seeding completed successfully!")
	return nil
}
package config

import (
	"fmt"
	"log"

	"dashboard-ac-backend/internal/domain"
	"dashboard-ac-backend/pkg/hash"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	Environment string `mapstructure:"ENVIRONMENT"`
	Port        string `mapstructure:"PORT"`
	JWTSecret   string `mapstructure:"JWT_SECRET"`
	Database    DatabaseConfig
}

type DatabaseConfig struct {
	Host     string `mapstructure:"DB_HOST"`
	Port     string `mapstructure:"DB_PORT"`
	User     string `mapstructure:"DB_USER"`
	Password string `mapstructure:"DB_PASSWORD"`
	Name     string `mapstructure:"DB_NAME"`
	SSLMode  string `mapstructure:"DB_SSLMODE"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	// Set default values
	viper.SetDefault("ENVIRONMENT", "development")
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("JWT_SECRET", "your-secret-key")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASSWORD", "password")
	viper.SetDefault("DB_NAME", "dashboard_ac")
	viper.SetDefault("DB_SSLMODE", "disable")

	// Try to read from environment-specific config file
	env := viper.GetString("ENVIRONMENT")
	viper.SetConfigName(fmt.Sprintf(".env.%s", env))
	viper.AddConfigPath(".")

	// Read config file if exists
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Config file not found, using environment variables and defaults: %v", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Unmarshal database config separately
	if err := viper.Unmarshal(&config.Database); err != nil {
		return nil, fmt.Errorf("failed to unmarshal database config: %w", err)
	}

	return &config, nil
}

func InitDatabase(cfg *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.Port,
		cfg.Database.SSLMode,
	)

	var gormLogger logger.Interface
	if cfg.Environment == "production" {
		gormLogger = logger.Default.LogMode(logger.Silent)
	} else {
		gormLogger = logger.Default.LogMode(logger.Info)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

// AutoMigrate runs database migrations automatically
func AutoMigrate(db *gorm.DB) error {
	log.Println("Starting database migration...")
	
	// Run auto migration for all models
	err := db.AutoMigrate(
		&domain.User{},
		&domain.Customer{},
		&domain.Technician{},
		&domain.Service{},
		&domain.Schedule{},
		&domain.Invoice{},
		&domain.InvoiceDetail{},
	)
	
	if err != nil {
		log.Printf("Failed to run migrations: %v", err)
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	
	log.Println("Database migration completed successfully")
	
	// Run seeding for initial data
	if err := seedInitialData(db); err != nil {
		log.Printf("Failed to seed initial data: %v", err)
		return fmt.Errorf("failed to seed initial data: %w", err)
	}
	
	return nil
}

// seedInitialData creates initial admin user if not exists
func seedInitialData(db *gorm.DB) error {
	log.Println("Checking for initial data...")
	
	// Check if admin user already exists
	var adminUser domain.User
	result := db.Where("email = ?", "admin@dashboardac.com").First(&adminUser)
	
	if result.Error == nil {
		log.Println("Admin user already exists, skipping seeding")
		return nil
	}
	
	// Create admin user
	hashedPassword, err := hash.HashPassword("admin123")
	if err != nil {
		return fmt.Errorf("failed to hash admin password: %w", err)
	}
	
	admin := domain.User{
		Name:     "Administrator",
		Email:    "admin@dashboardac.com",
		Password: hashedPassword,
		Role:     "admin",
	}
	
	if err := db.Create(&admin).Error; err != nil {
		return fmt.Errorf("failed to create admin user: %w", err)
	}
	
	log.Println("Initial admin user created successfully")
	log.Println("Admin credentials - Email: admin@dashboardac.com, Password: admin123")
	
	return nil
}
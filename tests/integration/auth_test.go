package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"dashboard-ac-backend/internal/api/handler"
	"dashboard-ac-backend/internal/api/middleware"
	"dashboard-ac-backend/internal/api/request"
	"dashboard-ac-backend/internal/api/response"
	"dashboard-ac-backend/internal/domain"
	"dashboard-ac-backend/internal/repository"
	"dashboard-ac-backend/internal/service"
	"dashboard-ac-backend/pkg/hash"
	"dashboard-ac-backend/pkg/logger"
)

type AuthTestSuite struct {
	suite.Suite
	app         *fiber.App
	db          *gorm.DB
	authHandler *handler.AuthHandler
}

func (suite *AuthTestSuite) SetupSuite() {
	// Set test environment
	os.Setenv("ENVIRONMENT", "test")
	os.Setenv("JWT_SECRET", "test-secret-key-for-integration-testing")

	// Initialize logger
	logger.InitLogger("test")

	// Setup in-memory SQLite database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	suite.Require().NoError(err)
	suite.db = db

	// Auto migrate
	err = db.AutoMigrate(&domain.User{})
	suite.Require().NoError(err)

	// Setup repositories and services
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, "test-secret-key-for-integration-testing")

	// Setup handlers
	suite.authHandler = handler.NewAuthHandler(authService)

	// Setup Fiber app
	suite.app = fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})

	// Setup routes
	api := suite.app.Group("/api/v1")
	auth := api.Group("/auth")
	auth.Post("/register", suite.authHandler.Register)
	auth.Post("/login", suite.authHandler.Login)
	auth.Post("/refresh", suite.authHandler.RefreshToken)
	auth.Get("/profile", middleware.JWTAuth("test-secret-key-for-integration-testing"), suite.authHandler.Me)
}

func (suite *AuthTestSuite) TearDownSuite() {
	sqlDB, _ := suite.db.DB()
	sqlDB.Close()
}

func (suite *AuthTestSuite) SetupTest() {
	// Clean database before each test
	suite.db.Exec("DELETE FROM users")
}

func (suite *AuthTestSuite) TestRegister_Success() {
	registerReq := request.RegisterRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
		Role:     "admin",
	}

	body, _ := json.Marshal(registerReq)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusCreated, resp.StatusCode)

	var response response.BaseResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	suite.NoError(err)
	suite.Equal("success", response.Status)
	suite.Equal("User registered successfully", response.Message)
}

func (suite *AuthTestSuite) TestRegister_DuplicateEmail() {
	// Create user first
	user := &domain.User{
		Name:     "Existing User",
		Email:    "test@example.com",
		Password: "hashedpassword",
		Role:     domain.RoleAdmin,
		IsActive: true,
	}
	suite.db.Create(user)

	registerReq := request.RegisterRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
		Role:     "admin",
	}

	body, _ := json.Marshal(registerReq)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusBadRequest, resp.StatusCode)
}

func (suite *AuthTestSuite) TestRegister_InvalidInput() {
	tests := []struct {
		name    string
		request request.RegisterRequest
	}{
		{
			name: "Empty name",
			request: request.RegisterRequest{
				Name:     "",
				Email:    "test@example.com",
				Password: "password123",
				Role:     "admin",
			},
		},
		{
			name: "Invalid email",
			request: request.RegisterRequest{
				Name:     "Test User",
				Email:    "invalid-email",
				Password: "password123",
				Role:     "admin",
			},
		},
		{
			name: "Short password",
			request: request.RegisterRequest{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "123",
				Role:     "admin",
			},
		},
		{
			name: "Invalid role",
			request: request.RegisterRequest{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "password123",
				Role:     "invalid",
			},
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			body, _ := json.Marshal(tt.request)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := suite.app.Test(req)
			suite.NoError(err)
			suite.Equal(http.StatusBadRequest, resp.StatusCode)
		})
	}
}

func (suite *AuthTestSuite) TestLogin_Success() {
	// Create user first
	hashedPassword, _ := hash.HashPassword("password123")
	user := &domain.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: hashedPassword,
		Role:     domain.RoleAdmin,
		IsActive: true,
	}
	suite.db.Create(user)

	loginReq := request.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	body, _ := json.Marshal(loginReq)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)

	var response struct {
		response.BaseResponse
		Data struct {
			User struct {
				ID    uint   `json:"id"`
				Name  string `json:"name"`
				Email string `json:"email"`
				Role  string `json:"role"`
			} `json:"user"`
			Tokens struct {
				AccessToken  string `json:"access_token"`
				RefreshToken string `json:"refresh_token"`
			} `json:"tokens"`
		} `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	suite.NoError(err)
	suite.Equal("success", response.Status)
	suite.NotEmpty(response.Data.Tokens.AccessToken)
	suite.NotEmpty(response.Data.Tokens.RefreshToken)
	suite.Equal("test@example.com", response.Data.User.Email)
}

func (suite *AuthTestSuite) TestLogin_InvalidCredentials() {
	// Create user first
	hashedPassword, _ := hash.HashPassword("password123")
	user := &domain.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: hashedPassword,
		Role:     domain.RoleAdmin,
		IsActive: true,
	}
	suite.db.Create(user)

	tests := []struct {
		name    string
		request request.LoginRequest
	}{
		{
			name: "Wrong password",
			request: request.LoginRequest{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
		},
		{
			name: "Non-existent email",
			request: request.LoginRequest{
				Email:    "nonexistent@example.com",
				Password: "password123",
			},
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			body, _ := json.Marshal(tt.request)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := suite.app.Test(req)
			suite.NoError(err)
			suite.Equal(http.StatusUnauthorized, resp.StatusCode)
		})
	}
}

func (suite *AuthTestSuite) TestGetProfile_Success() {
	// Create user first
	hashedPassword, _ := hash.HashPassword("password123")
	user := &domain.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: hashedPassword,
		Role:     domain.RoleAdmin,
		IsActive: true,
	}
	suite.db.Create(user)

	// Login to get token
	loginReq := request.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	body, _ := json.Marshal(loginReq)
	loginRequest := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
	loginRequest.Header.Set("Content-Type", "application/json")
	loginResp, _ := suite.app.Test(loginRequest)

	var loginResponse struct {
		Data struct {
			Tokens struct {
				AccessToken string `json:"access_token"`
			} `json:"tokens"`
		} `json:"data"`
	}
	json.NewDecoder(loginResp.Body).Decode(&loginResponse)

	// Test profile endpoint
	req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/profile", nil)
	req.Header.Set("Authorization", "Bearer "+loginResponse.Data.Tokens.AccessToken)

	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)

	var response struct {
		response.BaseResponse
		Data domain.User `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	suite.NoError(err)
	suite.Equal("success", response.Status)
	suite.Equal("test@example.com", response.Data.Email)
}

func (suite *AuthTestSuite) TestGetProfile_Unauthorized() {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/profile", nil)

	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(http.StatusUnauthorized, resp.StatusCode)
}

func TestAuthTestSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}
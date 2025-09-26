package service

import (
	"errors"

	"dashboard-ac-backend/internal/api/request"
	"dashboard-ac-backend/internal/domain"
	"dashboard-ac-backend/internal/repository"
	"dashboard-ac-backend/pkg/hash"
	"dashboard-ac-backend/pkg/jwt"

	"gorm.io/gorm"
)

type AuthService interface {
	Register(req *request.RegisterRequest) (*domain.User, error)
	Login(req *request.LoginRequest) (*jwt.TokenPair, *domain.User, error)
	RefreshToken(refreshToken string) (*jwt.TokenPair, error)
	ValidateToken(tokenString string) (*jwt.Claims, error)
}

type authService struct {
	userRepo  repository.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo repository.UserRepository, jwtSecret string) AuthService {
	return &authService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *authService) Register(req *request.RegisterRequest) (*domain.User, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Hash password
	hashedPassword, err := hash.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Set default role if not provided
	role := req.Role
	if role == "" {
		role = domain.RoleCustomer
	}

	// Create user
	user := &domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     role,
		IsActive: true,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(req *request.LoginRequest) (*jwt.TokenPair, *domain.User, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, errors.New("invalid email or password")
		}
		return nil, nil, err
	}

	// Check if user is active
	if !user.IsActive {
		return nil, nil, errors.New("user account is deactivated")
	}

	// Verify password
	if err := hash.CheckPassword(user.Password, req.Password); err != nil {
		return nil, nil, errors.New("invalid email or password")
	}

	// Generate token pair
	tokenPair, err := jwt.GenerateTokenPair(user, s.jwtSecret)
	if err != nil {
		return nil, nil, err
	}

	return tokenPair, user, nil
}

func (s *authService) RefreshToken(refreshToken string) (*jwt.TokenPair, error) {
	// Validate refresh token
	claims, err := jwt.ValidateToken(refreshToken, s.jwtSecret)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Check if it's a refresh token
	if claims.Subject != "refresh" {
		return nil, errors.New("invalid token type")
	}

	// Get user from database
	user, err := s.userRepo.GetByID(claims.UserID)
	if err != nil {
		return nil, err
	}

	// Check if user is still active
	if !user.IsActive {
		return nil, errors.New("user account is deactivated")
	}

	// Generate new token pair
	tokenPair, err := jwt.GenerateTokenPair(user, s.jwtSecret)
	if err != nil {
		return nil, err
	}

	return tokenPair, nil
}

func (s *authService) ValidateToken(tokenString string) (*jwt.Claims, error) {
	return jwt.ValidateToken(tokenString, s.jwtSecret)
}
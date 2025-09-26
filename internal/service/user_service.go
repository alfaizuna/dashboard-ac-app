package service

import (
	"errors"

	"dashboard-ac-backend/internal/api/request"
	"dashboard-ac-backend/internal/domain"
	"dashboard-ac-backend/internal/repository"
	"dashboard-ac-backend/pkg/hash"

	"gorm.io/gorm"
)

type UserService interface {
	Create(req *request.UserCreateRequest) (*domain.User, error)
	GetByID(id uint) (*domain.User, error)
	Update(id uint, req *request.UserUpdateRequest) (*domain.User, error)
	Delete(id uint) error
	List(pagination *request.PaginationRequest) ([]*domain.User, int64, error)
	GetByRole(role domain.Role, pagination *request.PaginationRequest) ([]*domain.User, int64, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) Create(req *request.UserCreateRequest) (*domain.User, error) {
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

	// Create user
	user := &domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     domain.Role(req.Role),
		IsActive: true,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) GetByID(id uint) (*domain.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return user, nil
}

func (s *userService) Update(id uint, req *request.UserUpdateRequest) (*domain.User, error) {
	// Get existing user
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// Update fields if provided
	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Email != nil {
		// Check if email is already taken by another user
		existingUser, err := s.userRepo.GetByEmail(*req.Email)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		if existingUser != nil && existingUser.ID != user.ID {
			return nil, errors.New("email is already taken")
		}
		user.Email = *req.Email
	}
	if req.Role != nil {
		user.Role = domain.Role(*req.Role)
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	// Save updated user
	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) Delete(id uint) error {
	// Check if user exists
	_, err := s.userRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	return s.userRepo.Delete(id)
}

func (s *userService) List(pagination *request.PaginationRequest) ([]*domain.User, int64, error) {
	offset := pagination.GetOffset()
	limit := pagination.GetLimit()
	
	return s.userRepo.List(offset, limit)
}

func (s *userService) GetByRole(role domain.Role, pagination *request.PaginationRequest) ([]*domain.User, int64, error) {
	offset := pagination.GetOffset()
	limit := pagination.GetLimit()
	
	return s.userRepo.GetByRole(role, offset, limit)
}
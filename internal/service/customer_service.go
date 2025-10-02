package service

import (
	"errors"

	"dashboard-ac-backend/internal/api/request"
	"dashboard-ac-backend/internal/domain"
	"dashboard-ac-backend/internal/repository"

	"gorm.io/gorm"
)

type CustomerService interface {
	Create(req *request.CustomerCreateRequest) (*domain.Customer, error)
	GetByID(id string) (*domain.Customer, error)
	Update(id string, req *request.CustomerUpdateRequest) (*domain.Customer, error)
	Delete(id string) error
	List(pagination *request.PaginationRequest) ([]*domain.Customer, int64, error)
	Search(req *request.CustomerSearchRequest) ([]*domain.Customer, int64, error)
}

type customerService struct {
	customerRepo repository.CustomerRepository
}

func NewCustomerService(customerRepo repository.CustomerRepository) CustomerService {
	return &customerService{
		customerRepo: customerRepo,
	}
}

func (s *customerService) Create(req *request.CustomerCreateRequest) (*domain.Customer, error) {
	customer := &domain.Customer{
		Name:    req.Name,
		Phone:   req.Phone,
		Address: req.Address,
		Email:   req.Email,
	}

	if err := s.customerRepo.Create(customer); err != nil {
		return nil, err
	}

	return customer, nil
}

func (s *customerService) GetByID(id string) (*domain.Customer, error) {
	customer, err := s.customerRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("customer not found")
		}
		return nil, err
	}

	return customer, nil
}

func (s *customerService) Update(id string, req *request.CustomerUpdateRequest) (*domain.Customer, error) {
	// Check if customer exists
	customer, err := s.customerRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("customer not found")
		}
		return nil, err
	}

	// Update fields if provided
	if req.Name != nil && *req.Name != "" {
		customer.Name = *req.Name
	}
	if req.Phone != nil && *req.Phone != "" {
		customer.Phone = *req.Phone
	}
	if req.Address != nil && *req.Address != "" {
		customer.Address = *req.Address
	}
	if req.Email != nil && *req.Email != "" {
		customer.Email = *req.Email
	}

	if err := s.customerRepo.Update(customer); err != nil {
		return nil, err
	}

	return customer, nil
}

func (s *customerService) Delete(id string) error {
	// Check if customer exists
	_, err := s.customerRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("customer not found")
		}
		return err
	}

	return s.customerRepo.Delete(id)
}

func (s *customerService) List(pagination *request.PaginationRequest) ([]*domain.Customer, int64, error) {
	offset := pagination.GetOffset()
	limit := pagination.GetLimit()

	return s.customerRepo.List(offset, limit)
}

func (s *customerService) Search(req *request.CustomerSearchRequest) ([]*domain.Customer, int64, error) {
	offset := req.GetOffset()
	limit := req.GetLimit()

	return s.customerRepo.Search(req.Name, req.Phone, req.Email, offset, limit)
}
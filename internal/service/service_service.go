package service

import (
	"errors"

	"dashboard-ac-backend/internal/api/request"
	"dashboard-ac-backend/internal/domain"
	"dashboard-ac-backend/internal/repository"

	"gorm.io/gorm"
)

type ServiceService interface {
	Create(req *request.ServiceCreateRequest) (*domain.Service, error)
	GetByID(id string) (*domain.Service, error)
	Update(id string, req *request.ServiceUpdateRequest) (*domain.Service, error)
	Delete(id string) error
	List(pagination *request.PaginationRequest) ([]*domain.Service, int64, error)
	Search(req *request.ServiceSearchRequest) ([]*domain.Service, int64, error)
}

type serviceService struct {
	serviceRepo repository.ServiceRepository
}

func NewServiceService(serviceRepo repository.ServiceRepository) ServiceService {
	return &serviceService{
		serviceRepo: serviceRepo,
	}
}

func (s *serviceService) Create(req *request.ServiceCreateRequest) (*domain.Service, error) {
	service := &domain.Service{
		Name:     req.Name,
		Price:    req.Price,
		Duration: req.Duration,
	}

	if err := s.serviceRepo.Create(service); err != nil {
		return nil, err
	}

	return service, nil
}

func (s *serviceService) GetByID(id string) (*domain.Service, error) {
	service, err := s.serviceRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("service not found")
		}
		return nil, err
	}

	return service, nil
}

func (s *serviceService) Update(id string, req *request.ServiceUpdateRequest) (*domain.Service, error) {
	// Check if service exists
	service, err := s.serviceRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("service not found")
		}
		return nil, err
	}

	// Update fields if provided
	if req.Name != nil && *req.Name != "" {
		service.Name = *req.Name
	}
	if req.Price != nil && *req.Price > 0 {
		service.Price = *req.Price
	}
	if req.Duration != nil && *req.Duration > 0 {
		service.Duration = *req.Duration
	}

	if err := s.serviceRepo.Update(service); err != nil {
		return nil, err
	}

	return service, nil
}

func (s *serviceService) Delete(id string) error {
	// Check if service exists
	_, err := s.serviceRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("service not found")
		}
		return err
	}

	return s.serviceRepo.Delete(id)
}

func (s *serviceService) List(pagination *request.PaginationRequest) ([]*domain.Service, int64, error) {
	offset := pagination.GetOffset()
	limit := pagination.GetLimit()

	return s.serviceRepo.List(offset, limit)
}

func (s *serviceService) Search(req *request.ServiceSearchRequest) ([]*domain.Service, int64, error) {
	offset := req.GetOffset()
	limit := req.GetLimit()

	return s.serviceRepo.Search(req.Name, req.MinPrice, req.MaxPrice, offset, limit)
}
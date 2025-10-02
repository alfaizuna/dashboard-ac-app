package service

import (
	"errors"

	"dashboard-ac-backend/internal/api/request"
	"dashboard-ac-backend/internal/domain"
	"dashboard-ac-backend/internal/repository"

	"gorm.io/gorm"
)

type TechnicianService interface {
	Create(req *request.TechnicianCreateRequest) (*domain.Technician, error)
	GetByID(id string) (*domain.Technician, error)
	Update(id string, req *request.TechnicianUpdateRequest) (*domain.Technician, error)
	Delete(id string) error
	List(pagination *request.PaginationRequest) ([]*domain.Technician, int64, error)
	Search(req *request.TechnicianSearchRequest) ([]*domain.Technician, int64, error)
}

type technicianService struct {
	technicianRepo repository.TechnicianRepository
}

func NewTechnicianService(technicianRepo repository.TechnicianRepository) TechnicianService {
	return &technicianService{
		technicianRepo: technicianRepo,
	}
}

func (s *technicianService) Create(req *request.TechnicianCreateRequest) (*domain.Technician, error) {
	technician := &domain.Technician{
		Name:           req.Name,
		Phone:          req.Phone,
		Specialization: req.Specialization,
	}

	if err := s.technicianRepo.Create(technician); err != nil {
		return nil, err
	}

	return technician, nil
}

func (s *technicianService) GetByID(id string) (*domain.Technician, error) {
	technician, err := s.technicianRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("technician not found")
		}
		return nil, err
	}

	return technician, nil
}

func (s *technicianService) Update(id string, req *request.TechnicianUpdateRequest) (*domain.Technician, error) {
	// Check if technician exists
	technician, err := s.technicianRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("technician not found")
		}
		return nil, err
	}

	// Update fields if provided
	if req.Name != nil && *req.Name != "" {
		technician.Name = *req.Name
	}
	if req.Phone != nil && *req.Phone != "" {
		technician.Phone = *req.Phone
	}
	if req.Specialization != nil && *req.Specialization != "" {
		technician.Specialization = *req.Specialization
	}

	if err := s.technicianRepo.Update(technician); err != nil {
		return nil, err
	}

	return technician, nil
}

func (s *technicianService) Delete(id string) error {
	// Check if technician exists
	_, err := s.technicianRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("technician not found")
		}
		return err
	}

	return s.technicianRepo.Delete(id)
}

func (s *technicianService) List(pagination *request.PaginationRequest) ([]*domain.Technician, int64, error) {
	offset := pagination.GetOffset()
	limit := pagination.GetLimit()

	return s.technicianRepo.List(offset, limit)
}

func (s *technicianService) Search(req *request.TechnicianSearchRequest) ([]*domain.Technician, int64, error) {
	offset := req.GetOffset()
	limit := req.GetLimit()

	return s.technicianRepo.Search(req.Name, req.Specialization, offset, limit)
}
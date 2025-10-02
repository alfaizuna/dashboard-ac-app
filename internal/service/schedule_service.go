package service

import (
	"errors"
	"time"

	"dashboard-ac-backend/internal/api/request"
	"dashboard-ac-backend/internal/domain"
	"dashboard-ac-backend/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ScheduleService interface {
	Create(req *request.ScheduleCreateRequest) (*domain.Schedule, error)
	GetByID(id string) (*domain.Schedule, error)
	Update(id string, req *request.ScheduleUpdateRequest) (*domain.Schedule, error)
	Delete(id string) error
	List(pagination *request.PaginationRequest) ([]*domain.Schedule, int64, error)
	Search(req *request.ScheduleSearchRequest) ([]*domain.Schedule, int64, error)
	GetByCustomerID(customerID string, pagination *request.PaginationRequest) ([]*domain.Schedule, int64, error)
	GetByTechnicianID(technicianID string, pagination *request.PaginationRequest) ([]*domain.Schedule, int64, error)
	GetByStatus(status domain.ScheduleStatus, pagination *request.PaginationRequest) ([]*domain.Schedule, int64, error)
}

type scheduleService struct {
	scheduleRepo   repository.ScheduleRepository
	customerRepo   repository.CustomerRepository
	technicianRepo repository.TechnicianRepository
	serviceRepo    repository.ServiceRepository
}

func NewScheduleService(
	scheduleRepo repository.ScheduleRepository,
	customerRepo repository.CustomerRepository,
	technicianRepo repository.TechnicianRepository,
	serviceRepo repository.ServiceRepository,
) ScheduleService {
	return &scheduleService{
		scheduleRepo:   scheduleRepo,
		customerRepo:   customerRepo,
		technicianRepo: technicianRepo,
		serviceRepo:    serviceRepo,
	}
}

func (s *scheduleService) Create(req *request.ScheduleCreateRequest) (*domain.Schedule, error) {
	// Parse UUIDs
	customerID, err := uuid.Parse(req.CustomerID)
	if err != nil {
		return nil, errors.New("invalid customer ID format")
	}

	technicianID, err := uuid.Parse(req.TechnicianID)
	if err != nil {
		return nil, errors.New("invalid technician ID format")
	}

	serviceID, err := uuid.Parse(req.ServiceID)
	if err != nil {
		return nil, errors.New("invalid service ID format")
	}

	// Validate customer exists
	_, err = s.customerRepo.GetByID(req.CustomerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("customer not found")
		}
		return nil, err
	}

	// Validate technician exists
	_, err = s.technicianRepo.GetByID(req.TechnicianID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("technician not found")
		}
		return nil, err
	}

	// Validate service exists
	_, err = s.serviceRepo.GetByID(req.ServiceID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("service not found")
		}
		return nil, err
	}

	schedule := &domain.Schedule{
		CustomerID:   customerID,
		TechnicianID: technicianID,
		ServiceID:    serviceID,
		Date:         req.Date,
		Time:         req.Time,
		Status:       domain.ScheduleStatusPending,
	}

	if err := s.scheduleRepo.Create(schedule); err != nil {
		return nil, err
	}

	return schedule, nil
}

func (s *scheduleService) GetByID(id string) (*domain.Schedule, error) {
	schedule, err := s.scheduleRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("schedule not found")
		}
		return nil, err
	}

	return schedule, nil
}

func (s *scheduleService) Update(id string, req *request.ScheduleUpdateRequest) (*domain.Schedule, error) {
	// Check if schedule exists
	schedule, err := s.scheduleRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("schedule not found")
		}
		return nil, err
	}

	// Update fields if provided
	if req.TechnicianID != nil && *req.TechnicianID != "" {
		technicianID, err := uuid.Parse(*req.TechnicianID)
		if err != nil {
			return nil, errors.New("invalid technician ID format")
		}
		
		// Validate technician exists
		_, err = s.technicianRepo.GetByID(*req.TechnicianID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("technician not found")
			}
			return nil, err
		}
		schedule.TechnicianID = technicianID
	}

	if req.ServiceID != nil && *req.ServiceID != "" {
		serviceID, err := uuid.Parse(*req.ServiceID)
		if err != nil {
			return nil, errors.New("invalid service ID format")
		}
		
		// Validate service exists
		_, err = s.serviceRepo.GetByID(*req.ServiceID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("service not found")
			}
			return nil, err
		}
		schedule.ServiceID = serviceID
	}

	if req.Date != nil && !req.Date.IsZero() {
		schedule.Date = *req.Date
	}

	if req.Time != nil && !req.Time.IsZero() {
		schedule.Time = *req.Time
	}

	if req.Status != nil && *req.Status != "" {
		status := domain.ScheduleStatus(*req.Status)
		// Validate status
		if status != domain.ScheduleStatusPending && 
		   status != domain.ScheduleStatusOnProgress && 
		   status != domain.ScheduleStatusCompleted && 
		   status != domain.ScheduleStatusCanceled {
			return nil, errors.New("invalid status")
		}
		schedule.Status = status
	}

	if err := s.scheduleRepo.Update(schedule); err != nil {
		return nil, err
	}

	return schedule, nil
}

func (s *scheduleService) Delete(id string) error {
	// Check if schedule exists
	_, err := s.scheduleRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("schedule not found")
		}
		return err
	}

	return s.scheduleRepo.Delete(id)
}

func (s *scheduleService) List(pagination *request.PaginationRequest) ([]*domain.Schedule, int64, error) {
	offset := pagination.GetOffset()
	limit := pagination.GetLimit()

	return s.scheduleRepo.List(offset, limit)
}

func (s *scheduleService) Search(req *request.ScheduleSearchRequest) ([]*domain.Schedule, int64, error) {
	offset := req.GetOffset()
	limit := req.GetLimit()

	var dateFrom, dateTo *time.Time
	if req.DateFrom != "" {
		if parsed, err := time.Parse("2006-01-02", req.DateFrom); err == nil {
			dateFrom = &parsed
		}
	}
	if req.DateTo != "" {
		if parsed, err := time.Parse("2006-01-02", req.DateTo); err == nil {
			dateTo = &parsed
		}
	}

	var status domain.ScheduleStatus
	if req.Status != "" {
		status = domain.ScheduleStatus(req.Status)
	}

	return s.scheduleRepo.Search(req.CustomerID, req.TechnicianID, req.ServiceID, status, dateFrom, dateTo, offset, limit)
}

func (s *scheduleService) GetByCustomerID(customerID string, pagination *request.PaginationRequest) ([]*domain.Schedule, int64, error) {
	offset := pagination.GetOffset()
	limit := pagination.GetLimit()
	
	return s.scheduleRepo.GetByCustomerID(customerID, offset, limit)
}

func (s *scheduleService) GetByTechnicianID(technicianID string, pagination *request.PaginationRequest) ([]*domain.Schedule, int64, error) {
	offset := pagination.GetOffset()
	limit := pagination.GetLimit()
	
	return s.scheduleRepo.GetByTechnicianID(technicianID, offset, limit)
}

func (s *scheduleService) GetByStatus(status domain.ScheduleStatus, pagination *request.PaginationRequest) ([]*domain.Schedule, int64, error) {
	offset := pagination.GetOffset()
	limit := pagination.GetLimit()
	
	return s.scheduleRepo.GetByStatus(status, offset, limit)
}
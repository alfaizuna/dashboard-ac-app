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

type InvoiceService interface {
	Create(req *request.InvoiceCreateRequest) (*domain.Invoice, error)
	GetByID(id string) (*domain.Invoice, error)
	Update(id string, req *request.InvoiceUpdateRequest) (*domain.Invoice, error)
	Delete(id string) error
	List(pagination *request.PaginationRequest) ([]*domain.Invoice, int64, error)
	Search(req *request.InvoiceSearchRequest) ([]*domain.Invoice, int64, error)
	GetByCustomerID(customerID string, pagination *request.PaginationRequest) ([]*domain.Invoice, int64, error)
	GetByScheduleID(scheduleID string) (*domain.Invoice, error)
	GetByStatus(status domain.InvoiceStatus, pagination *request.PaginationRequest) ([]*domain.Invoice, int64, error)
}

type invoiceService struct {
	invoiceRepo  repository.InvoiceRepository
	customerRepo repository.CustomerRepository
	scheduleRepo repository.ScheduleRepository
}

func NewInvoiceService(
	invoiceRepo repository.InvoiceRepository,
	customerRepo repository.CustomerRepository,
	scheduleRepo repository.ScheduleRepository,
) InvoiceService {
	return &invoiceService{
		invoiceRepo:  invoiceRepo,
		customerRepo: customerRepo,
		scheduleRepo: scheduleRepo,
	}
}

func (s *invoiceService) Create(req *request.InvoiceCreateRequest) (*domain.Invoice, error) {
	// Parse UUIDs
	customerID, err := uuid.Parse(req.CustomerID)
	if err != nil {
		return nil, errors.New("invalid customer ID format")
	}

	scheduleID, err := uuid.Parse(req.ScheduleID)
	if err != nil {
		return nil, errors.New("invalid schedule ID format")
	}

	// Validate customer exists
	_, err = s.customerRepo.GetByID(req.CustomerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("customer not found")
		}
		return nil, err
	}

	// Validate schedule exists
	_, err = s.scheduleRepo.GetByID(req.ScheduleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("schedule not found")
		}
		return nil, err
	}

	invoice := &domain.Invoice{
		CustomerID:  customerID,
		ScheduleID:  scheduleID,
		InvoiceDate: req.InvoiceDate,
		DueDate:     req.DueDate,
		Status:      domain.InvoiceStatusUnpaid,
		TotalAmount: 0, // Will be calculated when invoice details are added
	}

	if err := s.invoiceRepo.Create(invoice); err != nil {
		return nil, err
	}

	return invoice, nil
}

func (s *invoiceService) GetByID(id string) (*domain.Invoice, error) {
	invoice, err := s.invoiceRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invoice not found")
		}
		return nil, err
	}

	return invoice, nil
}

func (s *invoiceService) Update(id string, req *request.InvoiceUpdateRequest) (*domain.Invoice, error) {
	// Check if invoice exists
	invoice, err := s.invoiceRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invoice not found")
		}
		return nil, err
	}

	// Update fields if provided
	if req.InvoiceDate != nil && !req.InvoiceDate.IsZero() {
		invoice.InvoiceDate = *req.InvoiceDate
	}

	if req.DueDate != nil && !req.DueDate.IsZero() {
		invoice.DueDate = *req.DueDate
	}

	if req.Status != nil && *req.Status != "" {
		status := domain.InvoiceStatus(*req.Status)
		// Validate status
		if status != domain.InvoiceStatusUnpaid && 
		   status != domain.InvoiceStatusPaid && 
		   status != domain.InvoiceStatusOverdue {
			return nil, errors.New("invalid status")
		}
		invoice.Status = status
	}

	if err := s.invoiceRepo.Update(invoice); err != nil {
		return nil, err
	}

	return invoice, nil
}

func (s *invoiceService) Delete(id string) error {
	// Check if invoice exists
	_, err := s.invoiceRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("invoice not found")
		}
		return err
	}

	return s.invoiceRepo.Delete(id)
}

func (s *invoiceService) List(pagination *request.PaginationRequest) ([]*domain.Invoice, int64, error) {
	offset := pagination.GetOffset()
	limit := pagination.GetLimit()

	return s.invoiceRepo.List(offset, limit)
}

func (s *invoiceService) Search(req *request.InvoiceSearchRequest) ([]*domain.Invoice, int64, error) {
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

	var status domain.InvoiceStatus
	if req.Status != "" {
		status = domain.InvoiceStatus(req.Status)
	}

	return s.invoiceRepo.Search(req.CustomerID, req.ScheduleID, status, dateFrom, dateTo, offset, limit)
}

func (s *invoiceService) GetByCustomerID(customerID string, pagination *request.PaginationRequest) ([]*domain.Invoice, int64, error) {
	offset := pagination.GetOffset()
	limit := pagination.GetLimit()
	
	return s.invoiceRepo.GetByCustomerID(customerID, offset, limit)
}

func (s *invoiceService) GetByScheduleID(scheduleID string) (*domain.Invoice, error) {
	return s.invoiceRepo.GetByScheduleID(scheduleID)
}

func (s *invoiceService) GetByStatus(status domain.InvoiceStatus, pagination *request.PaginationRequest) ([]*domain.Invoice, int64, error) {
	offset := pagination.GetOffset()
	limit := pagination.GetLimit()
	
	return s.invoiceRepo.GetByStatus(status, offset, limit)
}
package service

import (
	"errors"

	"dashboard-ac-backend/internal/api/request"
	"dashboard-ac-backend/internal/domain"
	"dashboard-ac-backend/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InvoiceDetailService interface {
	Create(req *request.InvoiceDetailCreateRequest) (*domain.InvoiceDetail, error)
	GetByID(id string) (*domain.InvoiceDetail, error)
	Update(id string, req *request.InvoiceDetailUpdateRequest) (*domain.InvoiceDetail, error)
	Delete(id string) error
	GetByInvoiceID(invoiceID string) ([]*domain.InvoiceDetail, error)
	DeleteByInvoiceID(invoiceID string) error
}

type invoiceDetailService struct {
	invoiceDetailRepo repository.InvoiceDetailRepository
	invoiceRepo       repository.InvoiceRepository
	serviceRepo       repository.ServiceRepository
}

func NewInvoiceDetailService(
	invoiceDetailRepo repository.InvoiceDetailRepository,
	invoiceRepo repository.InvoiceRepository,
	serviceRepo repository.ServiceRepository,
) InvoiceDetailService {
	return &invoiceDetailService{
		invoiceDetailRepo: invoiceDetailRepo,
		invoiceRepo:       invoiceRepo,
		serviceRepo:       serviceRepo,
	}
}

func (s *invoiceDetailService) Create(req *request.InvoiceDetailCreateRequest) (*domain.InvoiceDetail, error) {
	// Parse UUIDs
	invoiceID, err := uuid.Parse(req.InvoiceID)
	if err != nil {
		return nil, errors.New("invalid invoice ID format")
	}

	serviceID, err := uuid.Parse(req.ServiceID)
	if err != nil {
		return nil, errors.New("invalid service ID format")
	}

	// Validate invoice exists
	_, err = s.invoiceRepo.GetByID(req.InvoiceID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invoice not found")
		}
		return nil, err
	}

	// Validate service exists
	service, err := s.serviceRepo.GetByID(req.ServiceID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("service not found")
		}
		return nil, err
	}

	// Calculate subtotal
	subtotal := service.Price * float64(req.Quantity)

	invoiceDetail := &domain.InvoiceDetail{
		InvoiceID: invoiceID,
		ServiceID: serviceID,
		Quantity:  req.Quantity,
		UnitPrice: service.Price,
		Subtotal:  subtotal,
	}

	if err := s.invoiceDetailRepo.Create(invoiceDetail); err != nil {
		return nil, err
	}

	// Update invoice total amount
	if err := s.updateInvoiceTotal(req.InvoiceID); err != nil {
		return nil, err
	}

	return invoiceDetail, nil
}

func (s *invoiceDetailService) GetByID(id string) (*domain.InvoiceDetail, error) {
	invoiceDetail, err := s.invoiceDetailRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invoice detail not found")
		}
		return nil, err
	}

	return invoiceDetail, nil
}

func (s *invoiceDetailService) Update(id string, req *request.InvoiceDetailUpdateRequest) (*domain.InvoiceDetail, error) {
	// Check if invoice detail exists
	invoiceDetail, err := s.invoiceDetailRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invoice detail not found")
		}
		return nil, err
	}

	// Update fields if provided
	if req.Quantity != nil && *req.Quantity > 0 {
		invoiceDetail.Quantity = *req.Quantity
		// Recalculate subtotal
		invoiceDetail.Subtotal = invoiceDetail.UnitPrice * float64(invoiceDetail.Quantity)
	}

	if req.UnitPrice != nil && *req.UnitPrice >= 0 {
		invoiceDetail.UnitPrice = *req.UnitPrice
		// Recalculate subtotal
		invoiceDetail.Subtotal = invoiceDetail.UnitPrice * float64(invoiceDetail.Quantity)
	}

	if err := s.invoiceDetailRepo.Update(invoiceDetail); err != nil {
		return nil, err
	}

	// Update invoice total amount
	if err := s.updateInvoiceTotal(invoiceDetail.InvoiceID.String()); err != nil {
		return nil, err
	}

	return invoiceDetail, nil
}

func (s *invoiceDetailService) Delete(id string) error {
	// Check if invoice detail exists
	invoiceDetail, err := s.invoiceDetailRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("invoice detail not found")
		}
		return err
	}

	invoiceID := invoiceDetail.InvoiceID.String()

	if err := s.invoiceDetailRepo.Delete(id); err != nil {
		return err
	}

	// Update invoice total amount
	if err := s.updateInvoiceTotal(invoiceID); err != nil {
		return err
	}

	return nil
}

func (s *invoiceDetailService) GetByInvoiceID(invoiceID string) ([]*domain.InvoiceDetail, error) {
	return s.invoiceDetailRepo.GetByInvoiceID(invoiceID)
}

func (s *invoiceDetailService) DeleteByInvoiceID(invoiceID string) error {
	return s.invoiceDetailRepo.DeleteByInvoiceID(invoiceID)
}

// Helper function to update invoice total amount
func (s *invoiceDetailService) updateInvoiceTotal(invoiceID string) error {
	// Get all invoice details for this invoice
	details, err := s.invoiceDetailRepo.GetByInvoiceID(invoiceID)
	if err != nil {
		return err
	}

	// Calculate total amount
	var totalAmount float64
	for _, detail := range details {
		totalAmount += detail.Subtotal
	}

	// Get invoice and update total amount
	invoice, err := s.invoiceRepo.GetByID(invoiceID)
	if err != nil {
		return err
	}

	invoice.TotalAmount = totalAmount
	return s.invoiceRepo.Update(invoice)
}
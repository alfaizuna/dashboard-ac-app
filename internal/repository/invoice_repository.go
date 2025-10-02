package repository

import (
	"time"

	"dashboard-ac-backend/internal/domain"

	"gorm.io/gorm"
)

type InvoiceRepository interface {
	Create(invoice *domain.Invoice) error
	GetByID(id string) (*domain.Invoice, error)
	Update(invoice *domain.Invoice) error
	Delete(id string) error
	List(offset, limit int) ([]*domain.Invoice, int64, error)
	Search(customerID, scheduleID string, status domain.InvoiceStatus, dateFrom, dateTo *time.Time, offset, limit int) ([]*domain.Invoice, int64, error)
	GetByCustomerID(customerID string, offset, limit int) ([]*domain.Invoice, int64, error)
	GetByScheduleID(scheduleID string) (*domain.Invoice, error)
	GetByStatus(status domain.InvoiceStatus, offset, limit int) ([]*domain.Invoice, int64, error)
}

type invoiceRepository struct {
	db *gorm.DB
}

func NewInvoiceRepository(db *gorm.DB) InvoiceRepository {
	return &invoiceRepository{db: db}
}

func (r *invoiceRepository) Create(invoice *domain.Invoice) error {
	return r.db.Create(invoice).Error
}

func (r *invoiceRepository) GetByID(id string) (*domain.Invoice, error) {
	var invoice domain.Invoice
	err := r.db.Where("id = ?", id).First(&invoice).Error
	if err != nil {
		return nil, err
	}
	return &invoice, nil
}

func (r *invoiceRepository) Update(invoice *domain.Invoice) error {
	return r.db.Save(invoice).Error
}

func (r *invoiceRepository) Delete(id string) error {
	return r.db.Delete(&domain.Invoice{}, "id = ?", id).Error
}

func (r *invoiceRepository) List(offset, limit int) ([]*domain.Invoice, int64, error) {
	var invoices []*domain.Invoice
	var total int64

	// Count total records
	if err := r.db.Model(&domain.Invoice{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated records
	err := r.db.Offset(offset).Limit(limit).Order("invoice_date DESC").Find(&invoices).Error
	if err != nil {
		return nil, 0, err
	}

	return invoices, total, nil
}

func (r *invoiceRepository) Search(customerID, scheduleID string, status domain.InvoiceStatus, dateFrom, dateTo *time.Time, offset, limit int) ([]*domain.Invoice, int64, error) {
	var invoices []*domain.Invoice
	var total int64

	query := r.db.Model(&domain.Invoice{})

	// Apply filters
	if customerID != "" {
		query = query.Where("customer_id = ?", customerID)
	}
	if scheduleID != "" {
		query = query.Where("schedule_id = ?", scheduleID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if dateFrom != nil {
		query = query.Where("invoice_date >= ?", dateFrom)
	}
	if dateTo != nil {
		query = query.Where("invoice_date <= ?", dateTo)
	}

	// Count total records with filters
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated records with filters
	err := query.Offset(offset).Limit(limit).Order("invoice_date DESC").Find(&invoices).Error
	if err != nil {
		return nil, 0, err
	}

	return invoices, total, nil
}

func (r *invoiceRepository) GetByCustomerID(customerID string, offset, limit int) ([]*domain.Invoice, int64, error) {
	var invoices []*domain.Invoice
	var total int64

	// Count total records
	if err := r.db.Model(&domain.Invoice{}).Where("customer_id = ?", customerID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated records
	err := r.db.Where("customer_id = ?", customerID).Offset(offset).Limit(limit).Order("invoice_date DESC").Find(&invoices).Error
	if err != nil {
		return nil, 0, err
	}

	return invoices, total, nil
}

func (r *invoiceRepository) GetByScheduleID(scheduleID string) (*domain.Invoice, error) {
	var invoice domain.Invoice
	err := r.db.Where("schedule_id = ?", scheduleID).First(&invoice).Error
	if err != nil {
		return nil, err
	}
	return &invoice, nil
}

func (r *invoiceRepository) GetByStatus(status domain.InvoiceStatus, offset, limit int) ([]*domain.Invoice, int64, error) {
	var invoices []*domain.Invoice
	var total int64

	// Count total records
	if err := r.db.Model(&domain.Invoice{}).Where("status = ?", status).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated records
	err := r.db.Where("status = ?", status).Offset(offset).Limit(limit).Order("invoice_date DESC").Find(&invoices).Error
	if err != nil {
		return nil, 0, err
	}

	return invoices, total, nil
}
package repository

import (
	"dashboard-ac-backend/internal/domain"

	"gorm.io/gorm"
)

type InvoiceDetailRepository interface {
	Create(invoiceDetail *domain.InvoiceDetail) error
	GetByID(id string) (*domain.InvoiceDetail, error)
	Update(invoiceDetail *domain.InvoiceDetail) error
	Delete(id string) error
	GetByInvoiceID(invoiceID string) ([]*domain.InvoiceDetail, error)
	DeleteByInvoiceID(invoiceID string) error
}

type invoiceDetailRepository struct {
	db *gorm.DB
}

func NewInvoiceDetailRepository(db *gorm.DB) InvoiceDetailRepository {
	return &invoiceDetailRepository{db: db}
}

func (r *invoiceDetailRepository) Create(invoiceDetail *domain.InvoiceDetail) error {
	return r.db.Create(invoiceDetail).Error
}

func (r *invoiceDetailRepository) GetByID(id string) (*domain.InvoiceDetail, error) {
	var invoiceDetail domain.InvoiceDetail
	err := r.db.Where("id = ?", id).First(&invoiceDetail).Error
	if err != nil {
		return nil, err
	}
	return &invoiceDetail, nil
}

func (r *invoiceDetailRepository) Update(invoiceDetail *domain.InvoiceDetail) error {
	return r.db.Save(invoiceDetail).Error
}

func (r *invoiceDetailRepository) Delete(id string) error {
	return r.db.Delete(&domain.InvoiceDetail{}, "id = ?", id).Error
}

func (r *invoiceDetailRepository) GetByInvoiceID(invoiceID string) ([]*domain.InvoiceDetail, error) {
	var invoiceDetails []*domain.InvoiceDetail
	err := r.db.Where("invoice_id = ?", invoiceID).Order("created_at ASC").Find(&invoiceDetails).Error
	if err != nil {
		return nil, err
	}
	return invoiceDetails, nil
}

func (r *invoiceDetailRepository) DeleteByInvoiceID(invoiceID string) error {
	return r.db.Delete(&domain.InvoiceDetail{}, "invoice_id = ?", invoiceID).Error
}
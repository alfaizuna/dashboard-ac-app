package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InvoiceDetail struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	InvoiceID uuid.UUID      `json:"invoice_id" gorm:"type:uuid;not null" validate:"required"`
	ServiceID uuid.UUID      `json:"service_id" gorm:"type:uuid;not null" validate:"required"`
	Quantity  int            `json:"quantity" gorm:"not null" validate:"required,min=1"`
	UnitPrice float64        `json:"unit_price" gorm:"type:decimal(10,2);not null" validate:"required,min=0"`
	Subtotal  float64        `json:"subtotal" gorm:"type:decimal(10,2);not null" validate:"required,min=0"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations will be handled at repository level to avoid circular imports
}

func (id *InvoiceDetail) BeforeCreate(tx *gorm.DB) error {
	if id.ID == uuid.Nil {
		id.ID = uuid.New()
	}
	// Calculate subtotal automatically
	id.Subtotal = float64(id.Quantity) * id.UnitPrice
	return nil
}

func (id *InvoiceDetail) BeforeUpdate(tx *gorm.DB) error {
	// Recalculate subtotal on update
	id.Subtotal = float64(id.Quantity) * id.UnitPrice
	return nil
}
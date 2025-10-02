package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InvoiceStatus string

const (
	InvoiceStatusUnpaid  InvoiceStatus = "Unpaid"
	InvoiceStatusPaid    InvoiceStatus = "Paid"
	InvoiceStatusOverdue InvoiceStatus = "Overdue"
)

type Invoice struct {
	ID           uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ScheduleID   uuid.UUID      `json:"schedule_id" gorm:"type:uuid;not null" validate:"required"`
	CustomerID   uuid.UUID      `json:"customer_id" gorm:"type:uuid;not null" validate:"required"`
	InvoiceDate  time.Time      `json:"invoice_date" gorm:"type:date;not null" validate:"required"`
	DueDate      time.Time      `json:"due_date" gorm:"type:date;not null" validate:"required"`
	TotalAmount  float64        `json:"total_amount" gorm:"type:decimal(10,2);not null" validate:"required,min=0"`
	Status       InvoiceStatus  `json:"status" gorm:"type:varchar(20);default:'Unpaid'" validate:"required,oneof=Unpaid Paid Overdue"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations will be handled at repository level to avoid circular imports
}

func (i *Invoice) BeforeCreate(tx *gorm.DB) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	if i.Status == "" {
		i.Status = InvoiceStatusUnpaid
	}
	return nil
}

func (i *Invoice) IsValidStatus() bool {
	return i.Status == InvoiceStatusUnpaid || 
		   i.Status == InvoiceStatusPaid || 
		   i.Status == InvoiceStatusOverdue
}
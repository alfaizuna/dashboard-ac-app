package request

import "time"

type InvoiceCreateRequest struct {
	ScheduleID  string    `json:"schedule_id" validate:"required,uuid"`
	CustomerID  string    `json:"customer_id" validate:"required,uuid"`
	InvoiceDate time.Time `json:"invoice_date" validate:"required"`
	DueDate     time.Time `json:"due_date" validate:"required"`
}

type InvoiceUpdateRequest struct {
	InvoiceDate *time.Time `json:"invoice_date,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	Status      *string    `json:"status,omitempty" validate:"omitempty,oneof=Unpaid Paid Overdue"`
}

type InvoiceSearchRequest struct {
	*PaginationRequest
	CustomerID  string `json:"customer_id" query:"customer_id"`
	ScheduleID  string `json:"schedule_id" query:"schedule_id"`
	Status      string `json:"status" query:"status"`
	DateFrom    string `json:"date_from" query:"date_from"`
	DateTo      string `json:"date_to" query:"date_to"`
}

type InvoiceDetailCreateRequest struct {
	InvoiceID string  `json:"invoice_id" validate:"required,uuid"`
	ServiceID string  `json:"service_id" validate:"required,uuid"`
	Quantity  int     `json:"quantity" validate:"required,min=1"`
	UnitPrice float64 `json:"unit_price" validate:"required,min=0"`
}

type InvoiceDetailUpdateRequest struct {
	Quantity  *int     `json:"quantity,omitempty" validate:"omitempty,min=1"`
	UnitPrice *float64 `json:"unit_price,omitempty" validate:"omitempty,min=0"`
}
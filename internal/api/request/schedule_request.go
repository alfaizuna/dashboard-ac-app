package request

import "time"

type ScheduleCreateRequest struct {
	CustomerID   string    `json:"customer_id" validate:"required,uuid"`
	TechnicianID string    `json:"technician_id" validate:"required,uuid"`
	ServiceID    string    `json:"service_id" validate:"required,uuid"`
	Date         time.Time `json:"date" validate:"required"`
	Time         time.Time `json:"time" validate:"required"`
}

type ScheduleUpdateRequest struct {
	TechnicianID *string    `json:"technician_id,omitempty" validate:"omitempty,uuid"`
	ServiceID    *string    `json:"service_id,omitempty" validate:"omitempty,uuid"`
	Date         *time.Time `json:"date,omitempty"`
	Time         *time.Time `json:"time,omitempty"`
	Status       *string    `json:"status,omitempty" validate:"omitempty,oneof=Pending On-Progress Completed Canceled"`
}

type ScheduleSearchRequest struct {
	*PaginationRequest
	CustomerID   string `json:"customer_id" query:"customer_id"`
	TechnicianID string `json:"technician_id" query:"technician_id"`
	ServiceID    string `json:"service_id" query:"service_id"`
	Status       string `json:"status" query:"status"`
	DateFrom     string `json:"date_from" query:"date_from"`
	DateTo       string `json:"date_to" query:"date_to"`
}
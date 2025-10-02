package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ScheduleStatus string

const (
	ScheduleStatusPending     ScheduleStatus = "Pending"
	ScheduleStatusOnProgress  ScheduleStatus = "On-Progress"
	ScheduleStatusCompleted   ScheduleStatus = "Completed"
	ScheduleStatusCanceled    ScheduleStatus = "Canceled"
)

type Schedule struct {
	ID           uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CustomerID   uuid.UUID      `json:"customer_id" gorm:"type:uuid;not null" validate:"required"`
	TechnicianID uuid.UUID      `json:"technician_id" gorm:"type:uuid;not null" validate:"required"`
	ServiceID    uuid.UUID      `json:"service_id" gorm:"type:uuid;not null" validate:"required"`
	Date         time.Time      `json:"date" gorm:"type:date;not null" validate:"required"`
	Time         time.Time      `json:"time" gorm:"type:time;not null" validate:"required"`
	Status       ScheduleStatus `json:"status" gorm:"type:varchar(20);default:'Pending'" validate:"required,oneof=Pending On-Progress Completed Canceled"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations will be handled at repository level to avoid circular imports
}

func (s *Schedule) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	if s.Status == "" {
		s.Status = ScheduleStatusPending
	}
	return nil
}

func (s *Schedule) IsValidStatus() bool {
	return s.Status == ScheduleStatusPending || 
		   s.Status == ScheduleStatusOnProgress || 
		   s.Status == ScheduleStatusCompleted || 
		   s.Status == ScheduleStatusCanceled
}
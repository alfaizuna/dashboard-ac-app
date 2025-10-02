package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name      string         `json:"name" gorm:"not null" validate:"required,min=2,max=100"`
	Price     float64        `json:"price" gorm:"type:decimal(10,2);not null" validate:"required,min=0"`
	Duration  int            `json:"duration" gorm:"not null" validate:"required,min=1"` // in minutes
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations will be handled at repository level to avoid circular imports
}

func (s *Service) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}
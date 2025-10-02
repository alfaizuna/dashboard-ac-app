package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Customer struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name      string         `json:"name" gorm:"not null" validate:"required,min=2,max=100"`
	Phone     string         `json:"phone" gorm:"uniqueIndex;not null" validate:"required,min=10,max=15"`
	Address   string         `json:"address" gorm:"type:text" validate:"required,min=10,max=500"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null" validate:"required,email"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations will be handled at repository level to avoid circular imports
}

func (c *Customer) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}
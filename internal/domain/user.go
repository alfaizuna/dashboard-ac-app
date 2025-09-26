package domain

import (
	"time"

	"gorm.io/gorm"
)

type Role string

const (
	RoleAdmin      Role = "admin"
	RoleTechnician Role = "technician"
	RoleCustomer   Role = "customer"
)

type User struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	Name      string         `json:"name" gorm:"not null"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null"`
	Password  string         `json:"-" gorm:"not null"`
	Role      Role           `json:"role" gorm:"type:varchar(20);default:'customer'"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (u *User) IsValidRole() bool {
	return u.Role == RoleAdmin || u.Role == RoleTechnician || u.Role == RoleCustomer
}

func (u *User) HasRole(role Role) bool {
	return u.Role == role
}

func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

func (u *User) IsTechnician() bool {
	return u.Role == RoleTechnician
}

func (u *User) IsCustomer() bool {
	return u.Role == RoleCustomer
}
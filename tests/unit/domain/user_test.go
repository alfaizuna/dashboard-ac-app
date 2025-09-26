package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"dashboard-ac-backend/internal/domain"
)

func TestUser_IsValidRole(t *testing.T) {
	tests := []struct {
		name     string
		role     domain.Role
		expected bool
	}{
		{
			name:     "Valid admin role",
			role:     domain.RoleAdmin,
			expected: true,
		},
		{
			name:     "Valid technician role",
			role:     domain.RoleTechnician,
			expected: true,
		},
		{
			name:     "Valid customer role",
			role:     domain.RoleCustomer,
			expected: true,
		},
		{
			name:     "Invalid role",
			role:     domain.Role("invalid"),
			expected: false,
		},
		{
			name:     "Empty role",
			role:     domain.Role(""),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &domain.User{Role: tt.role}
			result := user.IsValidRole()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUser_IsAdmin(t *testing.T) {
	tests := []struct {
		name     string
		role     domain.Role
		expected bool
	}{
		{
			name:     "Admin role",
			role:     domain.RoleAdmin,
			expected: true,
		},
		{
			name:     "Technician role",
			role:     domain.RoleTechnician,
			expected: false,
		},
		{
			name:     "Customer role",
			role:     domain.RoleCustomer,
			expected: false,
		},
		{
			name:     "Invalid role",
			role:     domain.Role("invalid"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &domain.User{Role: tt.role}
			result := user.IsAdmin()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUser_IsTechnician(t *testing.T) {
	tests := []struct {
		name     string
		role     domain.Role
		expected bool
	}{
		{
			name:     "Technician role",
			role:     domain.RoleTechnician,
			expected: true,
		},
		{
			name:     "Admin role",
			role:     domain.RoleAdmin,
			expected: false,
		},
		{
			name:     "Customer role",
			role:     domain.RoleCustomer,
			expected: false,
		},
		{
			name:     "Invalid role",
			role:     domain.Role("invalid"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &domain.User{Role: tt.role}
			result := user.IsTechnician()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUser_Creation(t *testing.T) {
	user := &domain.User{
		ID:        1,
		Name:      "Test User",
		Email:     "test@example.com",
		Password:  "hashedpassword",
		Role:      domain.RoleAdmin,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	assert.Equal(t, uint(1), user.ID)
	assert.Equal(t, "Test User", user.Name)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "hashedpassword", user.Password)
	assert.Equal(t, domain.RoleAdmin, user.Role)
	assert.True(t, user.IsActive)
	assert.True(t, user.IsAdmin())
	assert.False(t, user.IsTechnician())
	assert.True(t, user.IsValidRole())
}

func TestUser_SoftDelete(t *testing.T) {
	now := time.Now()
	user := &domain.User{
		ID:        1,
		Name:      "Test User",
		Email:     "test@example.com",
		Password:  "hashedpassword",
		Role:      domain.RoleAdmin,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: gorm.DeletedAt{Time: now, Valid: true},
	}

	assert.True(t, user.DeletedAt.Valid)
	assert.Equal(t, now, user.DeletedAt.Time)
}
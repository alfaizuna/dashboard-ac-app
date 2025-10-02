package repository

import (
	"dashboard-ac-backend/internal/domain"

	"gorm.io/gorm"
)

type TechnicianRepository interface {
	Create(technician *domain.Technician) error
	GetByID(id string) (*domain.Technician, error)
	Update(technician *domain.Technician) error
	Delete(id string) error
	List(offset, limit int) ([]*domain.Technician, int64, error)
	Search(name, specialization string, offset, limit int) ([]*domain.Technician, int64, error)
}

type technicianRepository struct {
	db *gorm.DB
}

func NewTechnicianRepository(db *gorm.DB) TechnicianRepository {
	return &technicianRepository{db: db}
}

func (r *technicianRepository) Create(technician *domain.Technician) error {
	return r.db.Create(technician).Error
}

func (r *technicianRepository) GetByID(id string) (*domain.Technician, error) {
	var technician domain.Technician
	err := r.db.Where("id = ?", id).First(&technician).Error
	if err != nil {
		return nil, err
	}
	return &technician, nil
}

func (r *technicianRepository) Update(technician *domain.Technician) error {
	return r.db.Save(technician).Error
}

func (r *technicianRepository) Delete(id string) error {
	return r.db.Delete(&domain.Technician{}, "id = ?", id).Error
}

func (r *technicianRepository) List(offset, limit int) ([]*domain.Technician, int64, error) {
	var technicians []*domain.Technician
	var total int64

	// Count total records
	if err := r.db.Model(&domain.Technician{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated records
	err := r.db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&technicians).Error
	if err != nil {
		return nil, 0, err
	}

	return technicians, total, nil
}

func (r *technicianRepository) Search(name, specialization string, offset, limit int) ([]*domain.Technician, int64, error) {
	var technicians []*domain.Technician
	var total int64

	query := r.db.Model(&domain.Technician{})

	// Apply filters
	if name != "" {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}
	if specialization != "" {
		query = query.Where("specialization ILIKE ?", "%"+specialization+"%")
	}

	// Count total records with filters
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated records with filters
	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&technicians).Error
	if err != nil {
		return nil, 0, err
	}

	return technicians, total, nil
}
package repository

import (
	"dashboard-ac-backend/internal/domain"

	"gorm.io/gorm"
)

type ServiceRepository interface {
	Create(service *domain.Service) error
	GetByID(id string) (*domain.Service, error)
	Update(service *domain.Service) error
	Delete(id string) error
	List(offset, limit int) ([]*domain.Service, int64, error)
	Search(name string, minPrice, maxPrice float64, offset, limit int) ([]*domain.Service, int64, error)
}

type serviceRepository struct {
	db *gorm.DB
}

func NewServiceRepository(db *gorm.DB) ServiceRepository {
	return &serviceRepository{db: db}
}

func (r *serviceRepository) Create(service *domain.Service) error {
	return r.db.Create(service).Error
}

func (r *serviceRepository) GetByID(id string) (*domain.Service, error) {
	var service domain.Service
	err := r.db.Where("id = ?", id).First(&service).Error
	if err != nil {
		return nil, err
	}
	return &service, nil
}

func (r *serviceRepository) Update(service *domain.Service) error {
	return r.db.Save(service).Error
}

func (r *serviceRepository) Delete(id string) error {
	return r.db.Delete(&domain.Service{}, "id = ?", id).Error
}

func (r *serviceRepository) List(offset, limit int) ([]*domain.Service, int64, error) {
	var services []*domain.Service
	var total int64

	// Count total records
	if err := r.db.Model(&domain.Service{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated records
	err := r.db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&services).Error
	if err != nil {
		return nil, 0, err
	}

	return services, total, nil
}

func (r *serviceRepository) Search(name string, minPrice, maxPrice float64, offset, limit int) ([]*domain.Service, int64, error) {
	var services []*domain.Service
	var total int64

	query := r.db.Model(&domain.Service{})

	// Apply filters
	if name != "" {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}
	if minPrice > 0 {
		query = query.Where("price >= ?", minPrice)
	}
	if maxPrice > 0 {
		query = query.Where("price <= ?", maxPrice)
	}

	// Count total records with filters
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated records with filters
	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&services).Error
	if err != nil {
		return nil, 0, err
	}

	return services, total, nil
}
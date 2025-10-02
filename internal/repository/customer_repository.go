package repository

import (
	"dashboard-ac-backend/internal/domain"

	"gorm.io/gorm"
)

type CustomerRepository interface {
	Create(customer *domain.Customer) error
	GetByID(id string) (*domain.Customer, error)
	Update(customer *domain.Customer) error
	Delete(id string) error
	List(offset, limit int) ([]*domain.Customer, int64, error)
	Search(name, phone, email string, offset, limit int) ([]*domain.Customer, int64, error)
}

type customerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return &customerRepository{db: db}
}

func (r *customerRepository) Create(customer *domain.Customer) error {
	return r.db.Create(customer).Error
}

func (r *customerRepository) GetByID(id string) (*domain.Customer, error) {
	var customer domain.Customer
	err := r.db.Where("id = ?", id).First(&customer).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *customerRepository) Update(customer *domain.Customer) error {
	return r.db.Save(customer).Error
}

func (r *customerRepository) Delete(id string) error {
	return r.db.Delete(&domain.Customer{}, "id = ?", id).Error
}

func (r *customerRepository) List(offset, limit int) ([]*domain.Customer, int64, error) {
	var customers []*domain.Customer
	var total int64

	// Count total records
	if err := r.db.Model(&domain.Customer{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated records
	err := r.db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&customers).Error
	if err != nil {
		return nil, 0, err
	}

	return customers, total, nil
}

func (r *customerRepository) Search(name, phone, email string, offset, limit int) ([]*domain.Customer, int64, error) {
	var customers []*domain.Customer
	var total int64

	query := r.db.Model(&domain.Customer{})

	// Apply filters
	if name != "" {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}
	if phone != "" {
		query = query.Where("phone ILIKE ?", "%"+phone+"%")
	}
	if email != "" {
		query = query.Where("email ILIKE ?", "%"+email+"%")
	}

	// Count total records with filters
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated records with filters
	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&customers).Error
	if err != nil {
		return nil, 0, err
	}

	return customers, total, nil
}
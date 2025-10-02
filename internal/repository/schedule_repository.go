package repository

import (
	"time"

	"dashboard-ac-backend/internal/domain"

	"gorm.io/gorm"
)

type ScheduleRepository interface {
	Create(schedule *domain.Schedule) error
	GetByID(id string) (*domain.Schedule, error)
	Update(schedule *domain.Schedule) error
	Delete(id string) error
	List(offset, limit int) ([]*domain.Schedule, int64, error)
	Search(customerID, technicianID, serviceID string, status domain.ScheduleStatus, dateFrom, dateTo *time.Time, offset, limit int) ([]*domain.Schedule, int64, error)
	GetByCustomerID(customerID string, offset, limit int) ([]*domain.Schedule, int64, error)
	GetByTechnicianID(technicianID string, offset, limit int) ([]*domain.Schedule, int64, error)
	GetByStatus(status domain.ScheduleStatus, offset, limit int) ([]*domain.Schedule, int64, error)
}

type scheduleRepository struct {
	db *gorm.DB
}

func NewScheduleRepository(db *gorm.DB) ScheduleRepository {
	return &scheduleRepository{db: db}
}

func (r *scheduleRepository) Create(schedule *domain.Schedule) error {
	return r.db.Create(schedule).Error
}

func (r *scheduleRepository) GetByID(id string) (*domain.Schedule, error) {
	var schedule domain.Schedule
	err := r.db.Where("id = ?", id).First(&schedule).Error
	if err != nil {
		return nil, err
	}
	return &schedule, nil
}

func (r *scheduleRepository) Update(schedule *domain.Schedule) error {
	return r.db.Save(schedule).Error
}

func (r *scheduleRepository) Delete(id string) error {
	return r.db.Delete(&domain.Schedule{}, "id = ?", id).Error
}

func (r *scheduleRepository) List(offset, limit int) ([]*domain.Schedule, int64, error) {
	var schedules []*domain.Schedule
	var total int64

	// Count total records
	if err := r.db.Model(&domain.Schedule{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated records
	err := r.db.Offset(offset).Limit(limit).Order("date DESC, time DESC").Find(&schedules).Error
	if err != nil {
		return nil, 0, err
	}

	return schedules, total, nil
}

func (r *scheduleRepository) Search(customerID, technicianID, serviceID string, status domain.ScheduleStatus, dateFrom, dateTo *time.Time, offset, limit int) ([]*domain.Schedule, int64, error) {
	var schedules []*domain.Schedule
	var total int64

	query := r.db.Model(&domain.Schedule{})

	// Apply filters
	if customerID != "" {
		query = query.Where("customer_id = ?", customerID)
	}
	if technicianID != "" {
		query = query.Where("technician_id = ?", technicianID)
	}
	if serviceID != "" {
		query = query.Where("service_id = ?", serviceID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if dateFrom != nil {
		query = query.Where("date >= ?", dateFrom)
	}
	if dateTo != nil {
		query = query.Where("date <= ?", dateTo)
	}

	// Count total records with filters
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated records with filters
	err := query.Offset(offset).Limit(limit).Order("date DESC, time DESC").Find(&schedules).Error
	if err != nil {
		return nil, 0, err
	}

	return schedules, total, nil
}

func (r *scheduleRepository) GetByCustomerID(customerID string, offset, limit int) ([]*domain.Schedule, int64, error) {
	var schedules []*domain.Schedule
	var total int64

	// Count total records
	if err := r.db.Model(&domain.Schedule{}).Where("customer_id = ?", customerID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated records
	err := r.db.Where("customer_id = ?", customerID).Offset(offset).Limit(limit).Order("date DESC, time DESC").Find(&schedules).Error
	if err != nil {
		return nil, 0, err
	}

	return schedules, total, nil
}

func (r *scheduleRepository) GetByTechnicianID(technicianID string, offset, limit int) ([]*domain.Schedule, int64, error) {
	var schedules []*domain.Schedule
	var total int64

	// Count total records
	if err := r.db.Model(&domain.Schedule{}).Where("technician_id = ?", technicianID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated records
	err := r.db.Where("technician_id = ?", technicianID).Offset(offset).Limit(limit).Order("date DESC, time DESC").Find(&schedules).Error
	if err != nil {
		return nil, 0, err
	}

	return schedules, total, nil
}

func (r *scheduleRepository) GetByStatus(status domain.ScheduleStatus, offset, limit int) ([]*domain.Schedule, int64, error) {
	var schedules []*domain.Schedule
	var total int64

	// Count total records
	if err := r.db.Model(&domain.Schedule{}).Where("status = ?", status).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated records
	err := r.db.Where("status = ?", status).Offset(offset).Limit(limit).Order("date DESC, time DESC").Find(&schedules).Error
	if err != nil {
		return nil, 0, err
	}

	return schedules, total, nil
}
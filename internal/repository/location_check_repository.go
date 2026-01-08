package database

import (
	"geowarns/internal/models"
	"gorm.io/gorm"
)

type LocationCheckRepository struct {
	db *gorm.DB
}

func NewLocationCheckRepository(db *gorm.DB) *LocationCheckRepository {
	return &LocationCheckRepository{db: db}
}

func (r *LocationCheckRepository) Create(check *models.LocationCheck) error {
	return r.db.Create(check).Error
}

func (r *LocationCheckRepository) GetAll() ([]models.LocationCheck, error) {
	var checks []models.LocationCheck
	err := r.db.Find(&checks).Error
	return checks, err
}

func (r *LocationCheckRepository) GetByID(id uint) (*models.LocationCheck, error) {
	var check models.LocationCheck
	err := r.db.First(&check, id).Error
	if err != nil {
		return nil, err
	}
	return &check, nil
}

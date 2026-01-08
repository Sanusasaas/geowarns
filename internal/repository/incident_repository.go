package database

import (
	"geowarns/internal/models"
	"gorm.io/gorm"
)

type IncidentRepository struct {
	db *gorm.DB
}

func NewIncidentRepository(db *gorm.DB) *IncidentRepository {
	return &IncidentRepository{db: db}
}

func (r *IncidentRepository) Create(incident *models.Incident) error {
	return r.db.Create(incident).Error
}

func (r *IncidentRepository) GetAll() ([]models.Incident, error) {
	var incidents []models.Incident
	err := r.db.Find(&incidents).Error
	return incidents, err
}

func (r *IncidentRepository) GetByID(id uint) (*models.Incident, error) {
	var incident models.Incident
	err := r.db.First(&incident, id).Error
	if err != nil {
		return nil, err
	}
	return &incident, nil
}

func (r *IncidentRepository) Update(incident *models.Incident) error {
	return r.db.Save(incident).Error
}

func (r *IncidentRepository) Delete(id uint) error {
	return r.db.Delete(&models.Incident{}, id).Error
}

// Метод для получения активных инцидентов
func (r *IncidentRepository) GetActiveIncidents() ([]models.Incident, error) {
	var incidents []models.Incident
	err := r.db.Where("is_active = ?", true).Find(&incidents).Error
	return incidents, err
}

package service

import (
	"geowarns/internal/models"
	repository "geowarns/internal/repository"
)

type LocationService struct {
	locationCheckRepo *repository.LocationCheckRepository
	incidentRepo      *repository.IncidentRepository
	webhookTaskRepo   *repository.WebhookTaskRepository
}

func NewLocationService(
	locationCheckRepo *repository.LocationCheckRepository,
	incidentRepo *repository.IncidentRepository,
	webhookTaskRepo *repository.WebhookTaskRepository,
) *LocationService {
	return &LocationService{
		locationCheckRepo: locationCheckRepo,
		incidentRepo:      incidentRepo,
		webhookTaskRepo:   webhookTaskRepo,
	}
}

func (s *LocationService) CheckLocation(req *models.LocationCheckRequest) (*models.LocationCheck, []models.Incident, error) {

	check := &models.LocationCheck{
		UserID:    req.UserID,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
	}

	if err := s.locationCheckRepo.Create(check); err != nil {
		return nil, nil, err
	}

	nearbyIncidents, err := s.findNearbyIncidents(req.Latitude, req.Longitude)
	if err != nil {
		return nil, nil, err
	}

	for _, incident := range nearbyIncidents {
		task := &models.WebhookTask{
			IncidentID: incident.ID,
			UserID:     req.UserID,
			Status:     "pending",
		}
		if err := s.webhookTaskRepo.Create(task); err != nil {
			continue
		}
	}

	return check, nearbyIncidents, nil
}

func (s *LocationService) findNearbyIncidents(lat, lng float64) ([]models.Incident, error) {
	return s.incidentRepo.GetActiveIncidents()
}

func (s *LocationService) GetLocationChecks() ([]models.LocationCheck, error) {
	return s.locationCheckRepo.GetAll()
}

func (s *LocationService) GetLocationCheckByID(id uint) (*models.LocationCheck, error) {
	return s.locationCheckRepo.GetByID(id)
}

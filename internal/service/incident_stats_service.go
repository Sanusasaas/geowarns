package service

import (
	"geowarns/internal/models"
	repository "geowarns/internal/repository"
)

type IncidentStatsService struct {
	statsRepo *repository.IncidentStatsRepository
}

func NewIncidentStatsService(statsRepo *repository.IncidentStatsRepository) *IncidentStatsService {
	return &IncidentStatsService{statsRepo: statsRepo}
}

func (s *IncidentStatsService) GetStats(timeWindowMinutes int) ([]models.IncidentStats, error) {
	return s.statsRepo.GetStats(timeWindowMinutes)
}

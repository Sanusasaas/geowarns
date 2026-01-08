package database

import (
	"fmt"
	"geowarns/internal/models"
	"time"

	"gorm.io/gorm"
)

type IncidentStatsRepository struct {
	db *gorm.DB
}

func NewIncidentStatsRepository(db *gorm.DB) *IncidentStatsRepository {
	return &IncidentStatsRepository{db: db}
}

func (r *IncidentStatsRepository) GetStats(timeWindowMinutes int) ([]models.IncidentStats, error) {
	now := time.Now()
	startTime := now.Add(-time.Duration(timeWindowMinutes) * time.Minute)

	type result struct {
		IncidentID  uint `gorm:"column:incident_id"`
		UserCount   int  `gorm:"column:user_count"`
	}

	var results []result
	query := fmt.Sprintf(`
		SELECT
			i.id as incident_id,
			COUNT(DISTINCT lc.user_id) as user_count
		FROM incidents i
		LEFT JOIN location_checks lc ON
			i.is_active = true AND
			lc.checked_at >= '%s' AND
			SQRT(POWER(lc.latitude - i.latitude, 2) + POWER(lc.longitude - i.longitude, 2)) <= i.radius/111.32
		GROUP BY i.id`,
		startTime.Format(time.RFC3339))

	if err := r.db.Raw(query).Scan(&results).Error; err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Преобразуем результаты в нужную структуру
	var stats []models.IncidentStats
	for _, r := range results {
		stats = append(stats, models.IncidentStats{
			IncidentID:  r.IncidentID,
			UserCount:   r.UserCount,
			TimeWindow:  timeWindowMinutes, // Устанавливаем правильное значение
			LastChecked: now.Format(time.RFC3339),
		})
	}

	return stats, nil
}

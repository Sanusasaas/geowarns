package database

import (
	"geowarns/internal/models"

	"gorm.io/gorm"
)

type WebhookTaskRepository struct {
	db *gorm.DB
}

func NewWebhookTaskRepository(db *gorm.DB) *WebhookTaskRepository {
	return &WebhookTaskRepository{db: db}
}

func (r *WebhookTaskRepository) Create(task *models.WebhookTask) error {
	return r.db.Create(task).Error
}

func (r *WebhookTaskRepository) GetPendingTasks() ([]models.WebhookTask, error) {
	var tasks []models.WebhookTask
	err := r.db.
		Where("status = ?", "pending").
		Order("created_at ASC").
		Limit(100).
		Find(&tasks).Error
	return tasks, err
}


func (r *WebhookTaskRepository) UpdateStatus(taskID uint, status string) error {
	return r.db.
		Model(&models.WebhookTask{}).
		Where("id = ?", taskID).
		Updates(map[string]interface{}{
			"status": status,
			"updated_at": gorm.Expr("NOW()"),
		}).Error
}

func (r *WebhookTaskRepository) GetIncidentByID(id uint) (*models.Incident, error) {
	var incident models.Incident
	if err := r.db.
		Where("id = ?", id).
		First(&incident).Error; err != nil {
		return nil, err
	}
	return &incident, nil
}

func (r *WebhookTaskRepository) GetTasksByStatus(status string, limit int) ([]models.WebhookTask, error) {
	var tasks []models.WebhookTask
	err := r.db.
		Where("status = ?", status).
		Order("created_at DESC").
		Limit(limit).
		Find(&tasks).Error
	return tasks, err
}

func (r *WebhookTaskRepository) GetTaskByID(id uint) (*models.WebhookTask, error) {
	var task models.WebhookTask
	if err := r.db.
		Where("id = ?", id).
		First(&task).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *WebhookTaskRepository) DeleteTask(id uint) error {
	return r.db.
		Where("id = ?", id).
		Delete(&models.WebhookTask{}).Error
}

func (r *WebhookTaskRepository) GetFailedTasks(limit int) ([]models.WebhookTask, error) {
	return r.GetTasksByStatus("failed", limit)
}

func (r *WebhookTaskRepository) GetCompletedTasks(limit int) ([]models.WebhookTask, error) {
	return r.GetTasksByStatus("completed", limit)
}

func (r *WebhookTaskRepository) GetStats() (map[string]int64, error) {
	var stats struct {
		Pending   int64
		Completed int64
		Failed    int64
	}

	err := r.db.
		Table("webhook_tasks").
		Select(`
			SUM(CASE WHEN status = 'pending' THEN 1 ELSE 0 END) as pending,
			SUM(CASE WHEN status = 'completed' THEN 1 ELSE 0 END) as completed,
			SUM(CASE WHEN status = 'failed' THEN 1 ELSE 0 END) as failed
		`).
		Row().
		Scan(&stats.Pending, &stats.Completed, &stats.Failed)

	if err != nil {
		return nil, err
	}

	return map[string]int64{
		"pending":   stats.Pending,
		"completed": stats.Completed,
		"failed":    stats.Failed,
	}, nil
}


// CleanOldTasks удаляет старые задачи ( > 30 дней)
func (r *WebhookTaskRepository) CleanOldTasks() error {
	return r.db.
		Where("created_at < NOW() - INTERVAL '30 days'").
		Delete(&models.WebhookTask{}).Error
}

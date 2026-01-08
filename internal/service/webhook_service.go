package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"geowarns/internal/models"
	repository "geowarns/internal/repository"

	"go.uber.org/zap"
)

type WebhookService struct {
	repo       *repository.WebhookTaskRepository
	httpClient *http.Client
	webhookURL string
	logger     *zap.Logger
}

func NewWebhookService(
	repo *repository.WebhookTaskRepository,
	webhookURL string,
	logger *zap.Logger,
) *WebhookService {
	return &WebhookService{
		repo:       repo,
		httpClient: &http.Client{Timeout: 30 * time.Second},
		webhookURL: webhookURL,
		logger:     logger,
	}
}

func (s *WebhookService) ProcessPendingTasks() error {
	tasks, err := s.repo.GetPendingTasks()
	if err != nil {
		return fmt.Errorf("failed to get pending tasks: %w", err)
	}

	for _, task := range tasks {
		go func(t models.WebhookTask) {
			if err := s.sendWebhook(t); err != nil {
				s.logger.Error("Failed to send webhook",
					zap.Uint("task_id", t.ID),
					zap.Error(err))
				if err := s.repo.UpdateStatus(t.ID, "failed"); err != nil {
					s.logger.Error("Failed to update task status",
						zap.Uint("task_id", t.ID),
						zap.Error(err))
				}
				return
			}

			if err := s.repo.UpdateStatus(t.ID, "completed"); err != nil {
				s.logger.Error("Failed to update task status",
					zap.Uint("task_id", t.ID),
					zap.Error(err))
			}
		}(task)
	}

	return nil
}

func (s *WebhookService) sendWebhook(task models.WebhookTask) error {
	incident, err := s.repo.GetIncidentByID(task.IncidentID)
	if err != nil {
		return fmt.Errorf("failed to get incident: %w", err)
	}

	payload := map[string]interface{}{
		"event":     "user_near_incident",
		"incident":  incident,
		"user_id":   task.UserID,
		"timestamp": time.Now().Format(time.RFC3339),
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	webhookURL, err := url.Parse(s.webhookURL)
	if err != nil {
		return fmt.Errorf("invalid webhook URL: %w", err)
	}

	req, err := http.NewRequest("POST", webhookURL.String(), bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("webhook returned status: %d", resp.StatusCode)
	}

	return nil
}

package handlers

import (
	"geowarns/internal/models"
	database "geowarns/internal/repository"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type WebhookHandler struct {
	repo   *database.WebhookTaskRepository
	logger *zap.Logger
}

func NewWebhookHandler(repo *database.WebhookTaskRepository, logger *zap.Logger) *WebhookHandler {
	return &WebhookHandler{
		repo:   repo,
		logger: logger,
	}
}

func (h *WebhookHandler) ProcessWebhook(c *fiber.Ctx) error {
	h.logger.Info("Processing webhook request",
		zap.String("method", c.Method()),
		zap.String("path", c.Path()),
		zap.String("ip", c.IP()))

	if c.Method() != fiber.MethodPost {
		h.logger.Warn("Invalid method for webhook", zap.String("method", c.Method()))
		return c.Status(fiber.StatusMethodNotAllowed).JSON(fiber.Map{
			"error": "Method not allowed",
		})
	}

	var payload map[string]interface{}
	if err := c.BodyParser(&payload); err != nil {
		h.logger.Error("Failed to parse webhook payload", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid payload",
		})
	}

	h.logger.Info("Webhook payload received", zap.Any("payload", payload))

	task := models.WebhookTask{
		IncidentID:  uint(payload["incident_id"].(float64)),
		UserID:      payload["user_id"].(string),
		Status:      "pending",
		Payload:     models.JSON(payload),
		Attempts:    0,
		NextAttempt: time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}


	if err := h.repo.Create(&task); err != nil {
		h.logger.Error("Failed to create webhook task", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to process webhook",
		})
	}

	h.logger.Info("Webhook task created successfully",
		zap.Uint("task_id", task.ID),
		zap.Uint("incident_id", task.IncidentID),
		zap.String("user_id", task.UserID))

	return c.JSON(fiber.Map{
		"status":  "queued",
		"message": "Webhook task created",
		"data": fiber.Map{
			"task_id": task.ID,
		},
	})
}

// HealthCheck для проверки состояния сервиса
func (h *WebhookHandler) HealthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":    "ok",
		"timestamp": time.Now().Format(time.RFC3339),
		"message":   "Webhook service is healthy",
	})
}

package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type HealthHandler struct {
	db *gorm.DB
}

func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

func (h *HealthHandler) HealthCheck(c *fiber.Ctx) error {
	dbStatus := "ok"
	dbError := ""
	if h.db != nil {
		sqlDB, err := h.db.DB()
		if err != nil {
			dbStatus = "error"
			dbError = err.Error()
		} else {
			if err := sqlDB.Ping(); err != nil {
				dbStatus = "error"
				dbError = err.Error()
			}
		}
	}

	response := fiber.Map{
		"status":    "ok",
		"message":   "System is healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"version":   "1.0.0",
		"components": fiber.Map{
			"database": fiber.Map{
				"status": dbStatus,
				"error":  dbError,
			},
		},
	}

	if dbStatus != "ok" {
		response["status"] = "degraded"
		response["message"] = "System is degraded"
	}

	return c.JSON(response)
}

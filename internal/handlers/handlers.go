package handlers

import (
	"net/http"
	"strconv"

	"geowarns/internal/models"
	database "geowarns/internal/repository"
	"geowarns/internal/service"

	"github.com/gofiber/fiber/v2"
)

type LocalRepository struct {
	db              *database.Repository
	locationService *service.LocationService
	statsService    *service.IncidentStatsService
}

func NewLocalRepository(
	db *database.Repository,
	locationService *service.LocationService,
	statsService *service.IncidentStatsService,
) *LocalRepository {
	return &LocalRepository{
		db:              db,
		locationService: locationService,
		statsService:    statsService,
	}
}

func (r *LocalRepository) CreateIncident(c *fiber.Ctx) error {
	var req models.IncidentCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "can't parse request",
			"error":   err.Error(),
		})
	}

	incident := &models.Incident{
		Title:       req.Title,
		Description: req.Description,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		Radius:      req.Radius,
		IsActive:    true,
	}

	if req.IsActive != nil {
		incident.IsActive = *req.IsActive
	}

	if err := r.db.DB.Create(incident).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "can't create incident",
			"error":   err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "incident was created successfully",
		"data":    incident,
	})
}


func (r *LocalRepository) GetIncidentList(c *fiber.Ctx) error {
	var incidents []models.Incident
	if err := r.db.DB.Find(&incidents).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "can't get incidents",
		})
	}

	return c.JSON(fiber.Map{
		"message": "incidents list",
		"data":    incidents,
	})
}

func (r *LocalRepository) GetIncidentByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid ID",
		})
	}

	var incident models.Incident
	if err := r.db.DB.First(&incident, id).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"message": "incident not found",
		})
	}

	return c.JSON(fiber.Map{
		"message": "incident found",
		"data":    incident,
	})
}

func (r *LocalRepository) UpdateIncidentByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid ID",
		})
	}

	var req models.IncidentCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "can't parse request",
		})
	}

	var incident models.Incident
	if err := r.db.DB.First(&incident, id).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"message": "incident not found",
		})
	}

	if req.Title != "" {
		incident.Title = req.Title
	}
	if req.Description != nil {
		incident.Description = req.Description
	}
	if req.Latitude != 0 {
		incident.Latitude = req.Latitude
	}
	if req.Longitude != 0 {
		incident.Longitude = req.Longitude
	}
	if req.Radius != 0 {
		incident.Radius = req.Radius
	}
	if req.IsActive != nil {
		incident.IsActive = *req.IsActive
	}

	if err := r.db.DB.Save(&incident).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "can't update incident",
		})
	}

	return c.JSON(fiber.Map{
		"message": "incident updated successfully",
		"data":    incident,
	})
}

func (r *LocalRepository) DeleteIncidentByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid ID",
		})
	}

	var incident models.Incident
	if err := r.db.DB.First(&incident, id).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"message": "incident not found",
		})
	}

	if err := r.db.DB.Delete(&incident).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "can't delete incident",
		})
	}

	return c.JSON(fiber.Map{
		"message": "incident deleted successfully",
	})
}

func (r *LocalRepository) CheckLocation(c *fiber.Ctx) error {
	var req models.LocationCheckRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "can't parse request",
		})
	}

	if req.UserID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "user_id is required",
		})
	}

	check, incidents, err := r.locationService.CheckLocation(&req)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "can't check location",
		})
	}

	return c.JSON(fiber.Map{
		"message": "location checked successfully",
		"data": fiber.Map{
			"check":     check,
			"incidents": incidents,
		},
	})
}

func (r *LocalRepository) GetLocationChecks(c *fiber.Ctx) error {
	checks, err := r.locationService.GetLocationChecks()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "can't get location checks",
		})
	}

	return c.JSON(fiber.Map{
		"message": "location checks list",
		"data":    checks,
	})
}

func (r *LocalRepository) GetLocationCheckByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid ID",
		})
	}

	check, err := r.locationService.GetLocationCheckByID(uint(id))
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"message": "location check not found",
		})
	}

	return c.JSON(fiber.Map{
		"message": "location check found",
		"data":    check,
	})
}

func (r *LocalRepository) GetIncidentStats(c *fiber.Ctx) error {
    timeWindow := c.Query("time_window", "30")
    timeWindowMinutes, err := strconv.Atoi(timeWindow)
    if err != nil {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{
            "message": "invalid time_window parameter",
        })
    }

    stats, err := r.statsService.GetStats(timeWindowMinutes)
    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
            "message": "can't get incident stats",
            "error":   err.Error(),
        })
    }

    return c.JSON(fiber.Map{
        "message": "incident stats retrieved successfully",
        "data":    stats,
    })
}


func (r *LocalRepository) SetupRoutes(app *fiber.App) {
	// Эндпоинты для инцидентов
	incidentAPI := app.Group("/api/v1/incidents")
	incidentAPI.Get("/stats", r.GetIncidentStats)
	incidentAPI.Post("/", r.CreateIncident)
	incidentAPI.Get("/", r.GetIncidentList)
	incidentAPI.Get("/:id", r.GetIncidentByID)
	incidentAPI.Put("/:id", r.UpdateIncidentByID)
	incidentAPI.Delete("/:id", r.DeleteIncidentByID)

	// Эндпоинты для проверки локаций
	locationAPI := app.Group("/api/v1/location")
	locationAPI.Post("/check", r.CheckLocation)
	locationAPI.Get("/", r.GetLocationChecks)
	locationAPI.Get("/:id", r.GetLocationCheckByID)

}

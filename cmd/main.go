package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"geowarns/internal/handlers"
	repository "geowarns/internal/repository"
	"geowarns/internal/service"
	"geowarns/migrations"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func main() {
	zapLogger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}
	defer zapLogger.Sync()

	if err := godotenv.Load(".env"); err != nil {
		zapLogger.Fatal("failed to load .env file", zap.Error(err))
	}

	dbConfig := repository.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	dbRepo, err := repository.NewConnection(&dbConfig)
	if err != nil {
		zapLogger.Fatal("failed to connect to database", zap.Error(err))
	}
	defer dbRepo.Close()

	zapLogger.Info("Running database migrations...")
	migrationFunctions := []struct {
		name string
		fn   func(*gorm.DB) error
	}{
		{"incidents", migrations.MigrateIncident},
		{"location_checks", migrations.MigrateLocationCheck},
		{"incident_stats", migrations.MigrateIncidentStat},
		{"webhook_tasks", migrations.MigrateWebhookTask},
	}

	var migrationErrs []error
	for _, mf := range migrationFunctions {
		if err := mf.fn(dbRepo.DB); err != nil {
			zapLogger.Error("Migration failed", zap.String("migration", mf.name), zap.Error(err))
			migrationErrs = append(migrationErrs, err)
		} else {
			zapLogger.Info("Migration completed", zap.String("migration", mf.name))
		}
	}

	if len(migrationErrs) > 0 {
		zapLogger.Fatal("One or more migrations failed", zap.Errors("errors", migrationErrs))
	}

	// Репозитории
	incidentRepo := repository.NewIncidentRepository(dbRepo.DB)
	locationCheckRepo := repository.NewLocationCheckRepository(dbRepo.DB)
	webhookTaskRepo := repository.NewWebhookTaskRepository(dbRepo.DB)
	incidentStatsRepo := repository.NewIncidentStatsRepository(dbRepo.DB)

	webhookURL := os.Getenv("WEBHOOK_URL")
	if webhookURL == "" {
		zapLogger.Warn("WEBHOOK_URL not set, using default", zap.String("url", "http://localhost:9090/webhook"))
		webhookURL = "http://localhost:9090/webhook"
	}

	// Сервисы
	statsService := service.NewIncidentStatsService(incidentStatsRepo)
	webhookService := service.NewWebhookService(webhookTaskRepo, webhookURL, zapLogger)
	locationService := service.NewLocationService(
		locationCheckRepo,
		incidentRepo,
		webhookTaskRepo,
	)

	// Хендлеры
	healthHandler := handlers.NewHealthHandler(dbRepo.DB)
	webhookHandler := handlers.NewWebhookHandler(webhookTaskRepo, zapLogger)
	mainHandler := handlers.NewLocalRepository(dbRepo, locationService, statsService)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				zapLogger.Info("Shutting down webhook processor...")
				return
			case <-ticker.C:
				if err := webhookService.ProcessPendingTasks(); err != nil {
					zapLogger.Error("Failed to process webhook tasks", zap.Error(err))
				}
			}
		}
	}()

	app := fiber.New(fiber.Config{
		AppName:      "GeoWarns API v1.0",
		BodyLimit:    10 * 1024 * 1024, // 10MB
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	})

	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path}\n",
	}))
	app.Use(func(c *fiber.Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.Next()
	})

	app.Get("/api/v1/system/health", healthHandler.HealthCheck)
	app.Post("/api/v1/webhooks", webhookHandler.ProcessWebhook)
	app.Get("/api/v1/webhooks/health", webhookHandler.HealthCheck)
	mainHandler.SetupRoutes(app)

	serverAddr := os.Getenv("SERVER_ADDR")
	if serverAddr == "" {
		serverAddr = ":8080"
	}

	zapLogger.Info("Starting server", zap.String("address", serverAddr))

	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-shutdownChan
		zapLogger.Info("Shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := app.ShutdownWithContext(ctx); err != nil {
			zapLogger.Error("Server shutdown error", zap.Error(err))
		}

		cancel()
		zapLogger.Info("Server shutdown complete")
	}()

	if err := app.Listen(serverAddr); err != nil {
		zapLogger.Error("Server error", zap.Error(err))
	}
}

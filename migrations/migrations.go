package migrations

import (
	"embed"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

//go:embed *.sql
var migrationsFS embed.FS

func MigrateIncident(db *gorm.DB) error {
	return runMigration(db, "01_incidents.sql")
}

func MigrateLocationCheck(db *gorm.DB) error {
	return runMigration(db, "02_location_checks.sql")
}

func MigrateIncidentStat(db *gorm.DB) error {
	return runMigration(db, "03_incident_stats.sql")
}

func MigrateWebhookTask(db *gorm.DB) error {
	return runMigration(db, "04_webhook_tasks.sql")
}

func runMigration(db *gorm.DB, filename string) error {
	data, err := migrationsFS.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	queries := strings.Split(string(data), ";")

	for _, query := range queries {
		query = strings.TrimSpace(query)
		if query == "" {
			continue
		}

		if err := db.Exec(query).Error; err != nil {
			return fmt.Errorf("failed to execute query '%s': %w", query, err)
		}
	}

	return nil
}

func MigrateAll(db *gorm.DB) error {
	migrations := []struct {
		name string
		fn   func(*gorm.DB) error
	}{
		{"incidents", MigrateIncident},
		{"location_checks", MigrateLocationCheck},
		{"incident_stats", MigrateIncidentStat},
		{"webhook_tasks", MigrateWebhookTask},
	}

	var errs []error
	for _, m := range migrations {
		if err := m.fn(db); err != nil {
			errs = append(errs, fmt.Errorf("migration %s failed: %w", m.name, err))
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

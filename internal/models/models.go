package models

import "time"

type Incident struct {
	ID            uint           `gorm:"primary_key" json:"id"`
	Title         string         `gorm:"not null" json:"title"`
	Description   *string         `json:"description"`
	Latitude      float64        `gorm:"not null" json:"latitude"`
	Longitude     float64        `gorm:"not null" json:"longitude"`
	Radius        float64        `gorm:"not null" json:"radius"`
	IsActive      bool           `gorm:"not null;default:true" json:"is_active"`
	CreatedAt     time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	WebhookTasks  []WebhookTask  `gorm:"foreignKey:IncidentID" json:"-"`
	IncidentStats []IncidentStat `gorm:"foreignKey:IncidentID" json:"-"`
}

type IncidentStat struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	IncidentID  uint      `gorm:"not null" json:"incident_id"`
	Incident    Incident  `gorm:"foreignKey:IncidentID" json:"-"`
	UserCount   int       `gorm:"default:0" json:"user_count"`
	LastUpdated time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"last_updated"`
}

type IncidentCreateRequest struct {
	Title       string  `json:"title" validate:"required,min=3,max=255"`
	Description *string `json:"description"`
	Latitude    float64 `json:"latitude" validate:"required,min=-90,max=90"`
	Longitude   float64 `json:"longitude" validate:"required,min=-180,max=180"`
	Radius      float64 `json:"radius" validate:"required,min=1"`
	IsActive    *bool   `json:"is_active"`
}

type IncidentStats struct {
	IncidentID  uint   `json:"incident_id"`
	UserCount   int    `json:"user_count"`
	TimeWindow  int    `json:"time_window_minutes"`
	LastChecked string `json:"last_checked"`
}

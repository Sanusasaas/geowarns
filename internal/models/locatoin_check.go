package models

import "time"

type LocationCheck struct {
    ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
    UserID    string    `json:"user_id"`
    Latitude  float64   `gorm:"not null" json:"latitude"`
    Longitude float64   `gorm:"not null" json:"longitude"`
    CheckedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"checked_at"`
}

type LocationCheckRequest struct {
	UserID    string  `json:"user_id" validate:"required"`
	Latitude  float64 `json:"latitude" validate:"required,min=-90,max=90"`
	Longitude float64 `json:"longitude" validate:"required,min=-180,max=180"`
}

package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type WebhookTask struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	IncidentID  uint      `json:"incident_id"`
	UserID      string    `json:"user_id"`
	Status      string    `gorm:"type:string;default:'pending'" json:"status"`
	Payload     JSON      `gorm:"type:jsonb" json:"payload"`
	Attempts    int       `gorm:"default:0" json:"attempts"`
	NextAttempt time.Time `json:"next_attempt"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type JSON map[string]interface{}

func (j JSON) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j *JSON) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, j)
}

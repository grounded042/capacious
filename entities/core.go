package entities

import "time"

type Event struct {
	EventId     string    `gorm:"primary_key" sql:"DEFAULT:uuid_generate_v1mc()" json:"event_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

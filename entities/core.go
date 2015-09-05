package entities

import "time"

type Event struct {
	EventId     string
	Name        string
	Description string
	StartTime   time.Time
	EndTime     time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

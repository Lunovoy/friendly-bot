package models

import (
	"github.com/google/uuid"
)

type Reminder struct {
	ID                uuid.UUID `json:"id" db:"id"`
	MinutesUntilEvent int       `json:"minutes_until_event" db:"minutes_until_event"`
	IsActive          bool      `json:"is_active" db:"is_active"`
	UserID            uuid.UUID `json:"user_id" db:"user_id"`
	EventID           uuid.UUID `json:"event_id" db:"event_id"`
}

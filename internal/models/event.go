package models

import (
	"database/sql"

	"github.com/google/uuid"
)

type Event struct {
	ID          uuid.UUID    `json:"id" db:"id"`
	Title       string       `json:"title" db:"title" binding:"required"`
	Description string       `json:"description" db:"description"`
	StartDate   sql.NullTime `json:"start_date" db:"start_date"`
	EndDate     sql.NullTime `json:"end_date" db:"end_date"`
	Frequency   *string      `json:"frequency" db:"frequency"`
	IsActive    bool         `json:"is_active" db:"is_active"`
	UserID      uuid.UUID    `json:"user_id" db:"user_id"`
}

type FriendsEvents struct {
	ID       uuid.UUID `json:"id," db:"id"`
	FriendID uuid.UUID `json:"friend_id" db:"friend_id"`
	EventID  uuid.UUID `json:"event_id" db:"event_id"`
}

type EventWithFriends struct {
	Event   Event    `json:"event"`
	Friends []Friend `json:"friends"`
}

type EventWithFriendsAndReminders struct {
	Event     Event      `json:"event"`
	Friends   []Friend   `json:"friends"`
	Reminders []Reminder `json:"reminders"`
}

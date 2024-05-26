package models

import (
	"database/sql"

	"github.com/google/uuid"
)

type Friend struct {
	ID        uuid.UUID    `json:"id,omitempty" db:"id"`
	FirstName string       `json:"first_name" db:"first_name"`
	LastName  string       `json:"last_name" db:"last_name"`
	DOB       sql.NullTime `json:"dob" db:"dob"`
	ImageID   uuid.UUID    `json:"image_id" db:"image_id"`
	UserID    uuid.UUID    `json:"user_id" db:"user_id"`
}

type FriendWithWorkInfo struct {
	ID                  uuid.UUID    `json:"id,omitempty" db:"id"`
	FirstName           string       `json:"first_name" db:"first_name"`
	LastName            string       `json:"last_name" db:"last_name"`
	DOB                 sql.NullTime `json:"dob" db:"dob"`
	ImageID             uuid.UUID    `json:"image_id" db:"image_id"`
	UserID              uuid.UUID    `json:"user_id" db:"user_id"`
	Messenger           string       `json:"messenger" db:"messenger"`
	CommunicationMethod string       `json:"communication_method" db:"communication_method"`
}

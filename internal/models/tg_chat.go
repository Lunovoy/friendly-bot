package models

import "github.com/google/uuid"

type TgChat struct {
	ID     uuid.UUID `json:"id,omitempty" db:"id"`
	ChatID int64     `json:"chat_id" db:"chat_id"`
	UserID uuid.UUID `json:"user_id" db:"user_id"`
}

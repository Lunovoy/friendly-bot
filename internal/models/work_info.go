package models

import "github.com/google/uuid"

type WorkInfo struct {
	ID                  uuid.UUID `json:"id" db:"id"`
	Country             string    `json:"country" db:"country"`
	City                string    `json:"city" db:"city"`
	Company             string    `json:"company" db:"company"`
	Profession          string    `json:"profession" db:"profession"`
	Position            string    `json:"position" db:"position"`
	Messenger           string    `json:"messenger" db:"messenger"`
	CommunicationMethod string    `json:"communication_method" db:"communication_method"`
	Nationality         string    `json:"nationality" db:"nationality"`
	Resident            bool      `json:"resident" db:"resident"`
	Language            string    `json:"language" db:"language"`
	FriendID            uuid.UUID `json:"friend_id" db:"friend_id"`
}

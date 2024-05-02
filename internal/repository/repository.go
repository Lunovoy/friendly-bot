package repository

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type User interface {
	GetUserByTelegramUsername(tgUsername string) (*uuid.UUID, error)
}

type Event interface {
}

type Repository struct {
	User
	Event
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Event: NewEventPostgres(db),
		User:  NewUserPostgres(db),
	}
}

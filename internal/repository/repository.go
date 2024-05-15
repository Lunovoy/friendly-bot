package repository

import (
	"friendly-bot/internal/models"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type User interface {
	GetUserByTelegramUsername(tgUsername string) (*uuid.UUID, error)
}

type TgChat interface {
	CreateTgChat(chatID int64, userID uuid.UUID) (uuid.UUID, error)
	GetTgChatByUserID(userID uuid.UUID) (*models.TgChat, error)
	GetTgChatByID(tgChatID uuid.UUID) (*models.TgChat, error)
	UpdateTgChat(chatID int64, userID uuid.UUID) error
}

type Event interface {
	GetEvents(currentTime time.Time) ([]*models.EventWithFriendsAndReminders, error)
	UpdateActiveStatus(eventID, userID uuid.UUID) error
	UpdateStartAndEndDate(eventID, userID uuid.UUID, startDate, endDate time.Time) error
}

type Repository struct {
	User
	TgChat
	Event
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		User:   NewUserPostgres(db),
		TgChat: NewTgChatPostgres(db),
		Event:  NewEventPostgres(db),
	}
}

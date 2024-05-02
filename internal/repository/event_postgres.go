package repository

import (
	"fmt"
	"friendly-bot/internal/models"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type EventPostgres struct {
	db *sqlx.DB
}

func NewEventPostgres(db *sqlx.DB) *EventPostgres {
	return &EventPostgres{
		db: db,
	}
}

func (r *EventPostgres) GetEvents(currentTime time.Time) (*[]models.Event, error) {
	query := fmt.Sprintf(`SELECT e.*, f.*
							FROM %s e
							JOIN %s fe ON e.id = fe.event_id
							JOIN %s f ON fe.friend_id = fe.friend_id
							WHERE e.start_date <= $1 AND e.end_date > $1`, eventTable, friendsEventsTable, friendTable)

	var data *[]models.Event
	err := r.db.Select(data, query, currentTime)

	return data, err
}

func (r *EventPostgres) UpdateStartEventStatus(eventID, userID uuid.UUID) error {

	query := fmt.Sprintf("UPDATE %s SET start_notify_sent = true id=$1 AND user_id=$2", eventTable)

	_, err := r.db.Exec(query, eventID, userID)

	return err
}

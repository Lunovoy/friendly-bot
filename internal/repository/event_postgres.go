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

func (r *EventPostgres) GetEvents(currentTime time.Time) ([]*models.EventWithFriends, error) {

	queryEvents := fmt.Sprintf(`SELECT e.*
								FROM %s e
								WHERE e.start_notify_sent = false AND 
								e.start_date <= $1 AND e.end_date > $1`, eventTable)

	queryFriends := fmt.Sprintf(`SELECT f.*
								FROM %s e
								JOIN %s fe ON e.id = fe.event_id
								JOIN %s f ON fe.friend_id = f.id
								WHERE fe.event_id = $1`, eventTable, friendsEventsTable, friendTable)

	eventWithFriends := []*models.EventWithFriends{}
	var events []*models.Event
	var friends []models.Friend

	err := r.db.Select(&events, queryEvents, currentTime)
	if err != nil {
		return nil, err
	}
	fmt.Println("Time: ", currentTime, "Events: ")
	for _, v := range events {
		fmt.Println(*v)
	}

	friendsStmt, err := r.db.Preparex(queryFriends)
	if err != nil {
		return nil, err
	}
	defer friendsStmt.Close()

	for _, event := range events {
		friendsStmt.Select(&friends, event.ID)
		eventWithFriends = append(eventWithFriends, &models.EventWithFriends{Event: *event, Friends: friends})
	}

	return eventWithFriends, err
}

func (r *EventPostgres) UpdateStartEventStatus(eventID, userID uuid.UUID) error {

	query := fmt.Sprintf("UPDATE %s SET start_notify_sent = true WHERE id = $1 AND user_id = $2", eventTable)

	_, err := r.db.Exec(query, eventID, userID)

	return err
}

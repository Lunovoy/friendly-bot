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

func (r *EventPostgres) GetEvents(currentTime time.Time) ([]*models.EventWithFriendsAndReminders, error) {

	// queryEvents := fmt.Sprintf(`SELECT e.*
	// 							FROM %s e
	// 							WHERE e.is_active = true AND
	// 							e.start_date <= $1 AND e.end_date > $1`, eventTable)

	queryEvents := fmt.Sprintf(`SELECT e.*
                            FROM %s e
                            LEFT JOIN %s r ON e.id = r.event_id
                            WHERE e.is_active = true AND r.is_active = true AND
                                  e.start_date <= NOW() + (r.minutes_until_event * INTERVAL '1 minute') AND 
                                  e.end_date > NOW()`, eventTable, reminderTable)

	queryFriends := fmt.Sprintf(`SELECT f.*
								FROM %s e
								JOIN %s fe ON e.id = fe.event_id
								JOIN %s f ON fe.friend_id = f.id
								WHERE fe.event_id = $1`, eventTable, friendsEventsTable, friendTable)

	queryReminders := fmt.Sprintf("SELECT * FROM %s WHERE event_id = $1 AND user_id = $2", reminderTable)

	eventWithFriendsAndReminders := []*models.EventWithFriendsAndReminders{}
	var events []*models.Event
	var friends []models.Friend
	var reminders []models.Reminder

	err := r.db.Select(&events, queryEvents)
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

	remindersStmt, err := r.db.Preparex(queryReminders)
	if err != nil {
		return nil, err
	}
	defer remindersStmt.Close()

	for _, event := range events {
		friendsStmt.Select(&friends, event.ID)
		remindersStmt.Select(&reminders, event.ID, event.UserID)
		eventWithFriendsAndReminders = append(eventWithFriendsAndReminders, &models.EventWithFriendsAndReminders{Event: *event, Friends: friends, Reminders: reminders})
		reminders = nil
		friends = nil
	}

	return eventWithFriendsAndReminders, err
}

func (r *EventPostgres) UpdateActiveStatus(eventID, userID uuid.UUID) error {

	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var count int
	queryCheck := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE event_id = $1 AND is_active = true", reminderTable)
	err = tx.Get(&count, queryCheck, eventID)
	fmt.Println("\n!!!Количество!!!: ", count, err)
	if err != nil {
		return err
	}

	if count != 0 {
		return err
	}

	query := fmt.Sprintf("UPDATE %s SET is_active = false WHERE id=$1 AND user_id=$2", eventTable)

	_, err = tx.Exec(query, eventID, userID)
	if err != nil {
		return err
	}

	queryReminders := fmt.Sprintf("UPDATE %s SET is_active = true WHERE event_id=$1", reminderTable)

	_, err = tx.Exec(queryReminders, eventID)
	if err != nil {
		return err
	}

	err = tx.Commit()

	return err
}

func (r *EventPostgres) UpdateStartAndEndDate(eventID, userID uuid.UUID, startDate, endDate time.Time) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var count int

	queryCheck := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE event_id = $1 AND is_active = true", reminderTable)
	_ = tx.Get(&count, queryCheck, eventID)
	fmt.Println("\n!!!Количество!!!: ", count, err)
	// if err != nil {
	// 	return err
	// }

	if count != 0 {
		return err
	}
	fmt.Println("Обновляем время!")
	query := fmt.Sprintf("UPDATE %s SET start_date = $1, end_date = $2 WHERE id=$3 AND user_id=$4", eventTable)

	_, err = tx.Exec(query, startDate, endDate, eventID, userID)
	if err != nil {
		return err
	}

	queryReminders := fmt.Sprintf("UPDATE %s SET is_active = true WHERE event_id=$1", reminderTable)

	_, err = tx.Exec(queryReminders, eventID)
	if err != nil {
		return err
	}

	err = tx.Commit()

	return err
}

func (r *EventPostgres) UpdateReminderStatus(reminderID uuid.UUID) error {

	query := fmt.Sprintf("UPDATE %s SET is_active = false WHERE id=$1", reminderTable)

	_, err := r.db.Exec(query, reminderID)

	return err
}

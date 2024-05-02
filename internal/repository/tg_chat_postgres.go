package repository

import (
	"fmt"
	"friendly-bot/internal/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type TgChatPostgres struct {
	db *sqlx.DB
}

func NewTgChatPostgres(db *sqlx.DB) *TgChatPostgres {
	return &TgChatPostgres{
		db: db,
	}
}

func (r *TgChatPostgres) CreateTgChat(chatID int64, userID uuid.UUID) (uuid.UUID, error) {

	query := fmt.Sprintf("INSERT INTO %s (chat_id, user_id) VALUES ($1, $2) RETURNING id", tgChatTable)

	var tgChatID uuid.UUID

	row := r.db.QueryRow(query, chatID, userID)
	if err := row.Scan(&tgChatID); err != nil {
		return uuid.Nil, err
	}

	return tgChatID, nil
}

func (r *TgChatPostgres) UpdateTgChat(chatID int64, userID uuid.UUID) error {

	query := fmt.Sprintf("UPDATE %s SET chat_id = $1 WHERE user_id = $2", tgChatTable)

	_, err := r.db.Exec(query, chatID, userID)

	return err
}

func (r *TgChatPostgres) GetTgChatByUserID(userID uuid.UUID) (*models.TgChat, error) {

	var tgChat models.TgChat
	query := fmt.Sprintf("SELECT * FROM \"%s\" WHERE user_id = $1", tgChatTable)

	err := r.db.Get(&tgChat, query, userID)
	return &tgChat, err
}

func (r *TgChatPostgres) GetTgChatByID(tgChatID uuid.UUID) (*models.TgChat, error) {

	var tgChat *models.TgChat
	query := fmt.Sprintf("SELECT * FROM \"%s\" WHERE id = $1", tgChatTable)

	err := r.db.Get(&tgChat, query, tgChatID)
	return tgChat, err
}

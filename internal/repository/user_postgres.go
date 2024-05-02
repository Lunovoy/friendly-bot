package repository

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserPostgres struct {
	db *sqlx.DB
}

func NewUserPostgres(db *sqlx.DB) *UserPostgres {
	return &UserPostgres{
		db: db,
	}
}

func (r *UserPostgres) GetUserByTelegramUsername(tgUsername string) (*uuid.UUID, error) {
	query := fmt.Sprintf("SELECT id from \"%s\" WHERE tg_username = $1", userTable)

	var userID *uuid.UUID
	err := r.db.Get(&userID, query, tgUsername)

	return userID, err
}

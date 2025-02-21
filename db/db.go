package db

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"os"
)

type Reminder struct {
	ID        int          `db:"id"`
	UserID    int64        `db:"user_id"`
	Username  string       `db:"username"`
	Day       string       `db:"day"`
	Time      string       `db:"time"`
	Remind1h  sql.NullTime `db:"remind_1h"`
	Remind24h sql.NullTime `db:"remind_24h"`
}

var DB *sqlx.DB

func Init() (*sqlx.DB, error) {
	connection := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("HOST"), os.Getenv("PORT"), os.Getenv("USER"), os.Getenv("PASSWORD"),
		os.Getenv("DBNAME"), os.Getenv("SSLMODE"))

	db, err := sqlx.Connect("postgres", connection)
	if err != nil {
		return nil, fmt.Errorf("Error connecting to database: %w", err)
	}
	DB = db
	return db, err
}

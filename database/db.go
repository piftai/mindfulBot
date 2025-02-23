package database

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"os"
	"time"
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
	log.Printf("Succesful authorization in DB")
	return db, err
}

func SaveReminder(userID int64, username string, day string, selectTime string) error {

	_, err := DB.Exec(`INSERT INTO reminders (user_id, username, day, time, remind_1h, remind_24h) VALUES ($1, $2, $3, $4, $5, $6)`, userID, username, day, selectTime, time.Now(), time.Now())
	return err
}

package models

import (
	"time"
)

type Reminder struct {
	ID        int       `db:"id"`
	UserID    int64     `db:"user_id"`
	Username  string    `db:"username"`
	Day       string    `db:"day"`
	Time      string    `db:"time"`
	Remind1h  time.Time `db:"remind_1h"`
	Remind24h time.Time `db:"remind_24h"`
}

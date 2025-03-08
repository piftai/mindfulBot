package models

import "database/sql"

type Reminder struct {
	ID        int          `db:"id"`
	UserID    int64        `db:"user_id"`
	Username  string       `db:"username"`
	Day       string       `db:"day"`
	Time      string       `db:"time"`
	Remind1h  sql.NullTime `db:"remind_1h"`
	Remind24h sql.NullTime `db:"remind_24h"`
}

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
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	DB = db
	log.Printf("Succesful authorization in DB")
	return db, err
}

func SaveReminder(userID int64, username, day, timeStr string) error {
	// Парсим время консультации
	consultationTime, err := time.Parse("15:04", timeStr)
	if err != nil {
		return fmt.Errorf("ошибка парсинга времени: %w", err)
	}
	log.Printf("SaveReminder is working: !!! consultationTime: %v\n", consultationTime)
	// Вычисляем текущую дату и время
	now := time.Now()

	// Определяем день консультации
	var consultationDate time.Time
	switch day {
	case "пн":
		consultationDate = nextWeekday(now, time.Monday)
	case "вт":
		consultationDate = nextWeekday(now, time.Tuesday)
	case "ср":
		consultationDate = nextWeekday(now, time.Wednesday)
	case "чт":
		consultationDate = nextWeekday(now, time.Thursday)
	case "пт":
		consultationDate = nextWeekday(now, time.Friday)
	default:
		return fmt.Errorf("неверный день недели: %s", day)
	}

	// Устанавливаем дату и время консультации
	consultationDateTime := time.Date(
		consultationDate.Year(),
		consultationDate.Month(),
		consultationDate.Day(),
		consultationTime.Hour(),
		consultationTime.Minute(),
		0, 0, time.Local,
	)

	// Вычисляем время напоминаний
	remind24h := consultationDateTime.Add(-24 * time.Hour)
	remind1h := consultationDateTime.Add(-1 * time.Hour)

	// Сохраняем напоминание в базу данных
	_, err = DB.Exec(`
		INSERT INTO reminders (user_id, username, day, time, remind_1h, remind_24h)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, userID, username, day, timeStr, remind1h, remind24h)
	return err
}

// Функция для вычисления следующего указанного дня недели
func nextWeekday(now time.Time, weekday time.Weekday) time.Time {
	daysUntilWeekday := (weekday - now.Weekday() + 7) % 7
	return now.Add(time.Duration(daysUntilWeekday) * 24 * time.Hour)
}

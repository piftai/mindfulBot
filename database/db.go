package database

import (
	"database/sql"
	"fmt"
	"log"
	"mindfulBot/models"
	"os"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

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

func SaveReminder(userID int64, username, day, timeStr string, isAlways bool) error {
	consultationTime, err := time.Parse("15:04", timeStr) // parsing consultation time
	if err != nil {
		return fmt.Errorf("ошибка парсинга времени: %w", err)
	}
	log.Printf("SaveReminder is working: !!! consultationTime: %v\n", consultationTime)
	now := time.Now()

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

	consultationDateTime := time.Date(
		consultationDate.Year(),
		consultationDate.Month(),
		consultationDate.Day(),
		consultationTime.Hour(),
		consultationTime.Minute(),
		0, 0, time.FixedZone("MSK", 3*60*60),
	)

	remind24h := consultationDateTime.Add(-24 * time.Hour)
	remind1h := consultationDateTime.Add(-1 * time.Hour)

	_, err = DB.Exec(`
		INSERT INTO reminders (user_id, username, day, time, remind_1h, remind_24h, is_always)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, userID, username, day, timeStr, remind1h, remind24h, isAlways)
	return err
}

func DeleteReminder(username string, reminderID string) (int64, error) {
	var value sql.Result
	var err error
	if reminderID == "all" {
		value, err = DB.Exec(`
			DELETE FROM reminders 
			WHERE username = $1
		`, username)
	} else {
		IntReminderID, _ := strconv.Atoi(reminderID)
		value, err = DB.Exec(`
			DELETE FROM reminders 
			WHERE username = $1 AND id = $2
		`, username, IntReminderID)
	}
	if err != nil {
		return 0, err
	}

	countOfDeletedRows, _ := value.RowsAffected()
	return countOfDeletedRows, nil
}

func nextWeekday(now time.Time, weekday time.Weekday) time.Time {
	daysUntilWeekday := (weekday - now.Weekday() + 7) % 7
	return now.Add(time.Duration(daysUntilWeekday) * 24 * time.Hour)
}

// --------------

func IsAdmin(username string) bool {
	val, err := DB.Exec(` 
		SELECT * FROM admin
		WHERE username = $1
	`, username)
	if err != nil {
		log.Printf("handleAdmin error: %v", err)
		return false
	}
	rows, err := val.RowsAffected()
	if err != nil {
		log.Printf("handleAdmin error: %v", err)
		return false
	}
	if rows == 1 {
		return true
	}
	return false
}

// --------------

func SaveUser(userID int64, username string) {
	_, err := DB.Exec(`
	INSERT INTO users (user_id, username)
	VALUES($1, $2)
	ON CONFLICT (user_id) DO NOTHING;
	`, userID, username)
	if err != nil {
		log.Printf("SaveUser error: %v", err)
	}
}

func GetUser(username string) (userID int64) {
	var user models.User
	err := DB.QueryRow(`
	SELECT user_id FROM users
	WHERE username = $1
	`, username).Scan(&user.UserID)

	if err != nil {
		log.Printf("GetUser error:%v", err)
	}

	return user.UserID
}

func GetUserById(id int) (userID int64) {
	var user models.User
	err := DB.QueryRow(`
	SELECT user_id FROM users
	WHERE id = $1
	`, id).Scan(&user.UserID)

	if err != nil {
		log.Printf("GetUser error:%v", err)
	}

	return user.UserID
}

func GetUsers() ([]models.User, error) {
	// Создаем слайс для хранения пользователей
	var users []models.User

	// Выполняем запрос
	rows, err := DB.Query("SELECT id, user_id FROM users WHERE username = '' OR username IS NULL")
	if err != nil {
		log.Printf("GetUsers query error: %v", err)
		return nil, err
	}
	defer rows.Close() // Закрываем результаты запроса

	// Итерируемся по строкам результата
	for rows.Next() {
		var user models.User
		// Сканируем данные в структуру
		if err := rows.Scan(&user.Id, &user.UserID); err != nil {
			log.Printf("GetUsers scan error: %v", err)
			return nil, err
		}
		// Добавляем пользователя в слайс
		users = append(users, user)
	}

	// Проверяем ошибки после итерации
	if err := rows.Err(); err != nil {
		log.Printf("GetUsers rows iteration error: %v", err)
		return nil, err
	}

	// Возвращаем слайс пользователей
	return users, nil
}

package scheduler

import (
	"database/sql"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jmoiron/sqlx"
	"github.com/robfig/cron/v3"
	"log"
	"mindfulBot/database"
	"time"
)

const paylink = "https://www.tinkoff.ru/rm/r_eKPOyRWmNB.XnfPKWHfzr/ZqYFh89264"

func Init(bot *tgbotapi.BotAPI, db *sqlx.DB) {
	c := cron.New()
	_, err := c.AddFunc("@every 1m", func() { checkReminders(bot, db) })
	if err != nil {
		log.Printf("Error add func in cron: %v", err)
	}
	c.Start()
	log.Println("Cron was started successful!")
}

func getReminders(db *sqlx.DB) ([]database.Reminder, error) {
	var reminders []database.Reminder
	err := db.Select(&reminders, `
	SELECT id, user_id, username, day, time, remind_1h, remind_24h
	FROM reminders
	WHERE remind_1h <= NOW() OR remind_24h <= NOW()
`)
	if err != nil {
		log.Printf("Error get reminders: %v", err)
		return nil, err
	}
	return reminders, nil
}

func checkReminders(bot *tgbotapi.BotAPI, db *sqlx.DB) {
	reminders, err := getReminders(db)
	if err != nil {
		log.Printf("Error check reminders: %v", err)
		return
	}
	for _, reminder := range reminders {
		sendReminder(bot, db, reminder)
	}
}

func sendReminder(bot *tgbotapi.BotAPI, db *sqlx.DB, reminder database.Reminder) {
	msgText := fmt.Sprintf("Ваша консультация состоится в %v %v.\n Ссылка на оплату: %v", reminder.Day, reminder.Time, paylink)
	msg := tgbotapi.NewMessage(reminder.UserID, msgText)
	bot.Send(msg)
	isUpdated, err := updateReminder(db, reminder)
	if !isUpdated {
		log.Printf("Reminder ID:%v did not update, but send.", reminder.ID)
	}
	if err != nil {
		log.Printf("Reminder ID:%v Error: %v", reminder.ID, err)
	}
}

func updateReminder(db *sqlx.DB, reminder database.Reminder) (bool, error) {
	isUpdated := false
	if !reminder.Remind24h.Time.After(time.Now()) {
		newRemind24h := sql.NullTime{
			Time:  reminder.Remind24h.Time.Add(7 * time.Hour * 24),
			Valid: true,
		}
		reminder.Remind24h = newRemind24h
		isUpdated = true
	}
	if !reminder.Remind1h.Time.After(time.Now()) {
		newRemind1h := sql.NullTime{
			Time:  reminder.Remind1h.Time.Add(7 * time.Hour * 24),
			Valid: true,
		}
		reminder.Remind1h = newRemind1h
		isUpdated = true
	}
	_, err := db.Exec(`
	UPDATE reminders
	SET remind_24h = $1,
		remind_1h = $2
	WHERE id = $3
	`, reminder.Remind24h, reminder.Remind1h, reminder.ID)
	if err != nil {
		log.Printf("Ошибка при обновлении напоминания %d: %v", reminder.ID, err)
		return false, err
	}
	return isUpdated, nil
}

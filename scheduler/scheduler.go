package scheduler

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jmoiron/sqlx"
	"github.com/robfig/cron/v3"
	"log"
	"mindfulBot/database"
)

func Init(bot *tgbotapi.BotAPI, db *sqlx.DB) {
	c := cron.New()
	_, err := c.AddFunc("@every 1m", func() { checkReminders(bot, db) })
	if err != nil {
		log.Printf("Error add func in cron: %v", err)
	}
	c.Start()
	log.Println("Cron was started succesful!")
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

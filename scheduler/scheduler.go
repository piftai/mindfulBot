package scheduler

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jmoiron/sqlx"
	"github.com/robfig/cron/v3"
	"log"
	"mindfulBot/database"
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
}

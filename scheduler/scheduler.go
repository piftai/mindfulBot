package scheduler

import (
	"database/sql"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jmoiron/sqlx"
	"github.com/robfig/cron/v3"
	"log"
	"mindfulBot/models"
	"time"
)

const paylink = "https://www.tinkoff.ru/rm/r_eKPOyRWmNB.XnfPKWHfzr/ZqYFh89264" // link to pay. not a secret so

func Init(bot *tgbotapi.BotAPI, db *sqlx.DB) {
	c := cron.New(cron.WithLocation(time.FixedZone("MSK", 3*60*60)))
	_, err := c.AddFunc("@every 1m", func() { checkReminders(bot, db) })
	if err != nil {
		log.Printf("Error add func in cron: %v", err)
	}
	c.Start()
	log.Println("Cron was started successful!")
}

func getReminders(db *sqlx.DB) ([]models.Reminder, error) {
	var reminders []models.Reminder
	err := db.Select(&reminders, `
	SELECT id, user_id, username, day, time, remind_1h, remind_24h
	FROM reminders
	WHERE remind_1h <= NOW() OR remind_24h <= NOW()
`)
	if err != nil {
		log.Printf("Error get reminders: %v", err)
		return nil, err
	}
	log.Printf("getReminders is working, current reminders: %v", reminders)
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

func sendReminder(bot *tgbotapi.BotAPI, db *sqlx.DB, reminder models.Reminder) {
	msgText := fmt.Sprintf("У тебя запланирована сессия.\n📅 День недели: %v\n🕒 "+
		"Время: %v\n💳 Ссылка на оплату: %v\n\n"+
		"Если у тебя изменились планы, напиши специалисту"+
		" в личку заранее, чтобы обсудить перенос.\n\n"+
		"Выдели это время только для себя. Найди спокойное место,"+
		" завари вкусный чай или просто настройся на работу с собой.\n "+
		"До встречи!", reminder.Day, reminder.Time, paylink)
	msg := tgbotapi.NewMessage(reminder.UserID, msgText)
	isUpdated, err := updateReminder(db, reminder)
	if !isUpdated {
		log.Printf("Reminder ID:%v did not update, and did not send.\n\nerror is: %v", reminder.ID, err)
		return
	}
	if err != nil {
		log.Printf("Reminder ID:%v Error: %v", reminder.ID, err)
		return
	}
	bot.Send(msg)
}

func updateReminder(db *sqlx.DB, reminder models.Reminder) (bool, error) {
	isUpdated := false
	log.Printf("reminder.Remind24h.Time.Before(time.Now()) = %v", reminder.Remind24h.Time.Before(time.Now()))
	if reminder.Remind24h.Time.Before(time.Now()) {
		newRemind24h := sql.NullTime{
			Time:  reminder.Remind24h.Time.Add(7 * 24 * time.Hour),
			Valid: true,
		}
		reminder.Remind24h = newRemind24h
		isUpdated = true
	}
	log.Printf("reminder.Remind1h.Time.Before(time.Now()) = %v", reminder.Remind1h.Time.Before(time.Now()))
	if reminder.Remind1h.Time.Before(time.Now()) {
		newRemind1h := sql.NullTime{
			Time:  reminder.Remind1h.Time.Add(7 * 24 * time.Hour),
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

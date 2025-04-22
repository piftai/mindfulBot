package handlers

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jmoiron/sqlx"
	"log"
	"mindfulBot/database"
	"strings"
)

func Router(bot *tgbotapi.BotAPI, message *tgbotapi.Message, db *sqlx.DB) {
	listWords := strings.Fields(message.Text)
	if len(listWords) > 1 {
		switch listWords[0] {

		}
	}
	switch message.Text {
	case "/start":
		handleStart(bot, message)
	case "/note":
		handleNote(bot, message)
	case "adminPing":
		handleAdmin(bot, message, db)
	}
}

type Bot interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
}

func handleAdmin(bot Bot, msg *tgbotapi.Message, db *sqlx.DB) {
	// перенести в пакет database
	val, err := db.Exec(` 
		SELECT * FROM admins
		WHERE username = $1
	`, msg.From.UserName)
	if err != nil {
		log.Printf("handleAdmin error: %v", err)
		return
	}
	rows, err := val.RowsAffected()
	if err != nil {
		log.Printf("handleAdmin error: %v", err)
		return
	}
	if rows == 1 {
		bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Admin pong successful"))
	}
}

func handleRandomNote(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {

}

func handleStart(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	message := tgbotapi.NewMessage(msg.Chat.ID, "Добро пожаловать в осознанные напоминания 🪷\n\n"+
		"Ты в сервисе записи к практическому психологу Софии. Здесь можно выбрать удобное время для сессии и получать напоминания о встречах.\n\n"+
		"— Запишись на свободное время\n"+
		"— Получай автоматические напоминания заранее\n"+
		"— Если нужно перенести встречу, сообщи об этом заранее\n\n"+
		"Нажми на команду ниже, чтобы выбрать время сессий 👇🏻\n/note")
	if msgClon, err := bot.Send(message); err != nil {
		log.Printf("message: %v. does not send: %v", msgClon, err)
	}
}

func handleNote(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	// creating buttons to pick a day
	days := []string{"пн", "вт", "ср", "чт", "пт"}
	var buttons []tgbotapi.InlineKeyboardButton

	for _, day := range days {
		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(day, "day_"+day))
	}

	// creating keyboard with buttons. that is inline buttons, so we need to turn on this func in botfather
	keyboard := tgbotapi.NewInlineKeyboardMarkup(buttons)
	reply := tgbotapi.NewMessage(msg.Chat.ID, "Давай выберем день для напоминания:")
	reply.ReplyMarkup = keyboard

	bot.Send(reply)
}

func HandleCallbackQuery(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	data := update.CallbackQuery.Data

	if strings.HasPrefix(data, "day_") {
		handleDaySelection(bot, update)
	} else if strings.HasPrefix(data, "time_") {
		handleTimeSelection(bot, update)
	}
}

func handleDaySelection(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	// to understand how it works check tg documentation
	day := strings.TrimPrefix(update.CallbackQuery.Data, "day_")

	times := getAvailableTimes(day)

	// creating buttons to pick a time
	var buttons []tgbotapi.InlineKeyboardButton
	for _, time := range times {
		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(time, "time_"+day+"_"+time))
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(buttons)
	reply := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Выбери время:")
	reply.ReplyMarkup = keyboard

	bot.Send(reply)
}

func handleTimeSelection(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	// extract day and time
	data := strings.TrimPrefix(update.CallbackQuery.Data, "time_")
	parts := strings.Split(data, "_")
	day := parts[0]
	time := parts[1]

	// save reminder into database
	userID := update.CallbackQuery.From.ID
	username := update.CallbackQuery.From.UserName
	err := database.SaveReminder(userID, username, day, time)
	if err != nil {
		log.Printf("Error saving notification: %v", err)
		bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Ошибка сохранения напоминания."))
		return
	}

	msg := fmt.Sprintf("Напоминание на %s %s сохранено. Если хочешь добавить "+
		"еще один день в свой календарь, то пропиши /note", day, time)

	bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, msg))
}

func getAvailableTimes(day string) []string {
	// slots for each day
	slots := map[string][]string{
		"пн": {"10:00", "12:00", "18:00"},
		"вт": {"11:00", "14:00", "16:00"},
		"ср": {"09:30", "11:30", "18:00"},
		"чт": {"10:00", "18:00"},
		"пт": {"10:00"},
	}

	return slots[day]
}

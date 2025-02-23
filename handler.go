package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"mindfulBot/database"
	"strings"
)

func router(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	switch message.Text {
	case "/start":
		handleStart(bot, message)
	case "/note":
		handleNote(bot, message)
	}
}

func handleStart(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	message := tgbotapi.NewMessage(msg.Chat.ID, "Привет, я бот, который будет напоминать о твоих сессиях!\n"+
		"Для начала нам нужно сделать с тобой твое расписание: для этого напиши в чат команду '/note'\n"+
		"После этого выбери удобный тебе день и время в которое ты хочешь заниматься.")
	if msgClon, err := bot.Send(message); err != nil {
		log.Printf("message: %v. does not send: %v", msgClon, err)
	}
}

func handleNote(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	// Создаем кнопки для выбора дня
	days := []string{"пн", "вт", "ср", "чт", "пт"}
	var buttons []tgbotapi.InlineKeyboardButton

	for _, day := range days {
		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(day, "day_"+day))
	}

	// Создаем клавиатуру с кнопками
	keyboard := tgbotapi.NewInlineKeyboardMarkup(buttons)
	reply := tgbotapi.NewMessage(msg.Chat.ID, "Привет, давай выберем день для напоминания:")
	reply.ReplyMarkup = keyboard

	// Отправляем сообщение
	bot.Send(reply)
}

func handleCallbackQuery(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	data := update.CallbackQuery.Data

	if strings.HasPrefix(data, "day_") {
		handleDaySelection(bot, update)
	} else if strings.HasPrefix(data, "time_") {
		handleTimeSelection(bot, update)
	}
}

func handleDaySelection(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	// Извлекаем выбранный день
	day := strings.TrimPrefix(update.CallbackQuery.Data, "day_")

	// Получаем доступные временные слоты для этого дня
	times := getAvailableTimes(day)

	// Создаем кнопки для выбора времени
	var buttons []tgbotapi.InlineKeyboardButton
	for _, time := range times {
		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(time, "time_"+day+"_"+time))
	}

	// Создаем клавиатуру с кнопками
	keyboard := tgbotapi.NewInlineKeyboardMarkup(buttons)
	reply := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Выбери время:")
	reply.ReplyMarkup = keyboard

	// Отправляем сообщение
	bot.Send(reply)
}

func handleTimeSelection(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	// Извлекаем день и время
	data := strings.TrimPrefix(update.CallbackQuery.Data, "time_")
	parts := strings.Split(data, "_")
	day := parts[0]
	time := parts[1]

	// Сохраняем напоминание в базу данных
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
	// Отправляем подтверждение
	bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, msg))
}

func getAvailableTimes(day string) []string {
	// Пример: временные слоты для каждого дня
	slots := map[string][]string{
		"пн": {"10:00", "12:00", "18:00"},
		"вт": {"11:00", "14:00", "16:00"},
		"ср": {"09:30", "11:30", "18:00"},
		"чт": {"10:00", "18:00"},
		"пт": {"10:00"},
	}

	return slots[day]
}

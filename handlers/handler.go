package handlers

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"mindfulBot/database"
	"strings"
)

/*
Пакет handler обеспечивает обработку введенной команды/набора символов/слов.
Это условный транпспортный слой, за маршрутизацию отвечает роутер, туда попадают все сообщения,
и внутри роутера мы определяем какая функция обработает сообщение. Будем называть функции обрабатывающие
сообщения хендлерами и называть handle<Something>

Функции внутри этого пакета неимпортируемые вне пакета, кроме Router & HandleCallbackQuery

HandleCallbackQuery нужен для работы с Inline buttons Telegram.
Это отдельный роутер для Inline data -> он доступен вне пакета
*/

func Router(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	database.SaveUser(message.From.ID, message.From.UserName)

	listWords := strings.Fields(message.Text)
	if len(listWords) > 1 {
		switch strings.ToLower(listWords[0]) {
		case "set":
			handleSet(bot, message)
		case "delete":
			handleDelete(bot, message)
		}
	}
	switch message.Text {
	case "/start":
		handleStart(bot, message)
	case "/note":
		handleNote(bot, message)
	case "adminPing":
		handleAdmin(bot, message)
	}
}

type Bot interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
}

func handleAdmin(bot Bot, msg *tgbotapi.Message) {
	if database.IsAdmin(msg.From.UserName) {
		if _, err := bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Admin pong successful")); err != nil {
			log.Printf("handleAdmin error: %v", err)
		}
	} else {
		if _, err := bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Admin pong unsuccessful. Check your permission to this operation")); err != nil {
			log.Printf("handleAdmin error: %v", err)
		}
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
	err := database.SaveReminder(userID, username, day, time, true) // true - потому что такого вида(запись от клиента, а не от админа) запись подразумевает напоминания всегда
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

func handleSet(bot Bot, msg *tgbotapi.Message) { // for admins only
	listWords := strings.Fields(msg.Text)
	var isAlwaysNotification bool
	for i := 0; i < len(listWords); i++ {
		listWords[i] = strings.ToLower(listWords[i])
	}

	if len(listWords) >= 5 { // опциональный параметр //
		if listWords[4] == "once" {
			isAlwaysNotification = false
		} else {
			isAlwaysNotification = true
		}
	} else {
		isAlwaysNotification = false
	}

	if database.IsAdmin(msg.From.UserName) {
		err := database.SaveReminder(msg.From.ID, listWords[1], listWords[2], listWords[3], isAlwaysNotification) // todo user id не клиента, а админа. надо поправить
		// пояснение к магическим цифрам выше:
		// 1 - ник клиента которому нужно поставить напоминание
		// 2 - день недели для напоминания
		// 3 - время напоминания
		if err != nil {
			bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Не удалось поставить напоминание. Попробуйте снова"))
			return
		}
		bot.Send(tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("Уведомление для %v успешно поставлено на %v %v", listWords[1], listWords[2], listWords[3])))
	} else {
		bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "У вас нет прав администратора для этой операции."))
	}
}

func handleDelete(bot Bot, msg *tgbotapi.Message) {
	listWords := strings.Fields(msg.Text)
	for i := 0; i < len(listWords); i++ {
		listWords[i] = strings.ToLower(listWords[i])
	}

	var username, id string

	if len(listWords) <= 1 {
		bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Неправильный формат команды\nПопробуйте delete @username id_уведомления/all"))
		return
	} else if len(listWords) == 2 {
		username = listWords[1]
		id = "all"
	} else if len(listWords) == 3 {
		username = listWords[1]
		id = listWords[2]
	} else {
		bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Неправильный формат команды\nПопробуйте ``` delete @username id_уведомления/all ```"))
		return
	}

	result, err := database.DeleteReminder(username, id)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "К сожалению, у меня не получилось удалить записи клиента\nПопробуйте снова"))
		return
	}

	bot.Send(tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("Отлично, у клиента %v удалилось %v записей.", username, result)))
}

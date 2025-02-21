package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"mindfulBot/db"
	"mindfulBot/utils"
	"os"
)

func main() {
	utils.Env()
	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	if db, err := db.Init(); err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			msg := tgbotapi.MessageConfig{}
			if update.Message.Text == "/start" {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Привет")
				bot.Send(msg)
			} else if update.Message.Text == "/note" {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Привет! Выбери день, в который мы будем встречаться")
				keyboard := tgbotapi.NewInlineKeyboardButtonData("понедельник", "понедельник")
				_ = keyboard
				bot.Send(msg)

			}
			//msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			//msg.ReplyToMessageID = update.Message.MessageID

		}
	}
}

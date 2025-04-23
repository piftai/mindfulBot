package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"mindfulBot/database"
	"mindfulBot/handlers"
	"mindfulBot/scheduler"
	"mindfulBot/utils"
	"os"
	"time"
)

func main() {
	utils.Env()
	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("time is:", time.Now()) // for debug. remove later
	db, err := database.Init()
	if err != nil {
		log.Panic(err)
	}
	_ = db

	log.Printf("Authorized on account %s", bot.Self.UserName)
	scheduler.Init(bot, db)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			handlers.Router(bot, update.Message)
		} else if update.CallbackQuery != nil {
			handlers.HandleCallbackQuery(bot, update)
		}
	}
}

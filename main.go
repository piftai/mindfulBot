package main

import (
	"context"
	"log"
	"mindfulBot/handlers"
	"mindfulBot/teachsimpleclient"
	"mindfulBot/utils"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const configRefreshInterval = 5 * time.Minute

func main() {
	utils.Env()
	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Println(err)
	}
	log.Println("Bot initialized")
	// db, err := database.Init()
	// if err != nil {
	// 	log.Println(err)
	// }
	// _ = db
	log.Println("Database initialized")
	// log.Printf("Authorized on account %s", bot.Self.UserName)

	proteachClient := teachsimpleclient.NewClient(os.Getenv("PROTEACH_URL"), os.Getenv("BOT_TOKEN"))
	configCache := teachsimpleclient.NewCache(proteachClient)
	if err := configCache.Refresh(context.Background()); err != nil {
		log.Println(err)
	}
	log.Println("Config fetched from proteach")
	configCache.StartAutoRefresh(context.Background(), configRefreshInterval)

	h := handlers.New(configCache)
	// scheduler.Init(bot, db, configCache)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	//updates := bot.GetUpdatesChan(u)

	var updates tgbotapi.UpdatesChannel

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			h.Router(bot, update.Message)
		} else if update.CallbackQuery != nil {
			h.HandleCallbackQuery(bot, update)
		}
	}
}

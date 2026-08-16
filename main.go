package main

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"mindfulBot/database"
	"mindfulBot/handlers"
	"mindfulBot/proteachclient"
	"mindfulBot/scheduler"
	"mindfulBot/utils"
	"os"
	"time"
)

const configRefreshInterval = 5 * time.Minute

func main() {
	utils.Env()
	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}
	log.Println("Bot initialized")
	db, err := database.Init()
	if err != nil {
		log.Panic(err)
	}
	_ = db
	log.Println("Database initialized")
	log.Printf("Authorized on account %s", bot.Self.UserName)

	proteachClient := proteachclient.NewClient(os.Getenv("PROTEACH_URL"), os.Getenv("BOT_TOKEN"))
	configCache := proteachclient.NewCache(proteachClient)
	if err := configCache.Refresh(context.Background()); err != nil {
		log.Panic(err)
	}
	log.Println("Config fetched from proteach")
	configCache.StartAutoRefresh(context.Background(), configRefreshInterval)

	h := handlers.New(configCache)
	scheduler.Init(bot, db, configCache)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			h.Router(bot, update.Message)
		} else if update.CallbackQuery != nil {
			h.HandleCallbackQuery(bot, update)
		}
	}
}

package bots

import (
	"mindfulBot/models"
)

// Repository - интерфейс к таблице bots.
type Repository interface {
	// Add - добавляет бота в репозиторий.
	Add(bot models.Bot) error
	// List - выводит список ботов.
	List() ([]models.Bot, error)

	// Flush - очищает таблицу.
	Flush() error
}

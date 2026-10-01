package models

import "time"

type Bot struct {
	// ID - идентифкатор бота.
	ID int
	// StID - идентификатор бота в simpleteach.
	StID int
	// OwnerID - идентификатор владельца в simpleteach.
	OwnerID int
	// Name - имя бота.
	Name string
	// Token - телеграмм токен бота.
	Token string

	CreatedAt time.Time
}

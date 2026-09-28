package postgres

import "time"

type bot struct {
	// ID - идентификатор бота.
	ID int `db:"id"`
	// StID - идентификатор бота в simpleteach.
	StID int `db:"st_id"`
	// OwnerID - идентификатор владельца в simpleteach.
	OwnerID int `db:"owner_id"`
	// Name - имя бота.
	Name string `db:"name"`
	// Token - телеграмм токен бота.
	Token string `db:"telegram_token"`

	CreatedAt time.Time `db:"created_at"`
}

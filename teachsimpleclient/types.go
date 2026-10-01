package teachsimpleclient

import "time"

// BotResponse - ответ на получение бота из simpleteach.
type BotResponse struct {
	// ID - идентификатор бота в simpleteach.
	ID int `json:"id"`
	// OwnerID - идентификатор владельца в simpleteach.
	OwnerID int `json:"owner_id"`
	// Name - имя бота.
	Name string `json:"name"`
	// Token - телеграмм токен бота.
	Token string `json:"token"`

	CreatedAt time.Time `json:"created_at"`
}
